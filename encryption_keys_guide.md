# Encryption Key Management Guide

This guide explains how to generate, store, and use AES-256 encryption keys for PHI (Protected Health Information) in the vcita-ai-agent-gin project.

## Overview

The application uses AES-256-GCM authenticated encryption to protect all PHI data at rest. This includes:
- Patient message content
- Medical notes and medication information
- Client contact details
- Any other sensitive healthcare data

**CRITICAL**: The encryption key is the only way to decrypt stored data. Loss of the key means permanent loss of all encrypted PHI.

## Key Generation

### Method 1: Using the Provided Script

```bash
# Generate a new random 32-byte AES-256 key
go run ./scripts/genkey
```

This outputs:
```
Generated AES-256 encryption key (base64):
your-base64-encoded-key-here

Add to your environment:
  ENCRYPTION_KEY_B64=your-base64-encoded-key-here
```

### Method 2: Using OpenSSL

```bash
# Generate 32 random bytes, hex-encoded
openssl rand -hex 32
```

### Method 3: Using /dev/urandom

```bash
# Generate 32 random bytes, base64-encoded
head -c 32 /dev/urandom | base64
```

## Key Storage

### Environment Variables

**Production**: Store the key in a secure secrets manager:
- AWS Secrets Manager
- HashiCorp Vault
- Azure Key Vault
- GCP Secret Manager
- Kubernetes Secrets (encrypted)

**Development**: Use `.env` file (never commit to git):

```bash
# .env
PHI_ENCRYPTION_KEY=your-64-character-hex-string-here
```

### Key Format Requirements

- **Length**: Exactly 32 bytes (256 bits)
- **Encoding**: Hexadecimal (64 characters) for environment variable
- **Example**: `PHI_ENCRYPTION_KEY=0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef`

The application automatically:
1. Reads the hex string from `PHI_ENCRYPTION_KEY` env var
2. Decodes it to 32 bytes
3. Uses it to initialize the AES-256-GCM cipher

## Key Usage in Code

### Initialization

```go
// In main.go
cfg, err := config.Load()
// cfg.PHIEncryptionKey is []byte (32 bytes)

enc, err := crypto.NewEncryptor(cfg.PHIEncryptionKey)
// enc is ready to encrypt/decrypt
```

### Encryption

```go
// Encrypt sensitive data before storing
encrypted, err := enc.Encrypt("patient message content")
if err != nil {
    return fmt.Errorf("encryption failed: %w", err)
}
// Store encrypted in database
```

### Decryption

```go
// Decrypt when retrieving data
decrypted, err := enc.Decrypt(encryptedFromDB)
if err != nil {
    return fmt.Errorf("decryption failed: %w", err)
}
// Use decrypted data
```

### Error Handling

The crypto package provides clear error messages:
- **Encryption failures**: Usually nonce generation issues (rare)
- **Decryption failures**: Authentication tag mismatch (tampered data or wrong key)

```go
plaintext, err := enc.Decrypt(ciphertext)
if err != nil {
    // Log error but DO NOT expose details in user-facing messages
    log.Error("PHI decryption failed", zap.Error(err))
    return errors.New("data integrity check failed")
}
```

## Key Rotation

**WARNING**: Key rotation is complex and requires careful planning.

### Process

1. **Generate new key**
2. **Add new environment variable** (e.g., `PHI_ENCRYPTION_KEY_V2`)
3. **Deploy application** that can read both keys
4. **Re-encrypt all data** with new key (background job)
5. **Verify all data** decrypts with new key
6. **Remove old key** from environment
7. **Clean up** old encrypted data if needed

### Implementation Considerations

```go
// Hypothetical dual-key encryptor
type DualEncryptor struct {
    oldKey, newKey *crypto.Encryptor
}

func (d *DualEncryptor) Decrypt(ciphertext string) (string, error) {
    // Try new key first
    if plain, err := d.newKey.Decrypt(ciphertext); err == nil {
        return plain, nil
    }
    // Fall back to old key
    return d.oldKey.Decrypt(ciphertext)
}
```

## Security Best Practices

### Key Generation
- Use cryptographically secure random sources
- Generate keys on the target system when possible
- Never reuse keys across environments

### Key Storage
- Never store keys in source code
- Never log keys or include in error messages
- Use dedicated secrets management systems
- Rotate keys regularly (quarterly minimum)

### Key Usage
- Initialize encryptor once at startup
- Use the same key for all operations during runtime
- Handle decryption errors gracefully without exposing key details
- Audit all encryption/decryption operations

### Operational Security
- Limit key access to necessary personnel
- Use HSMs or KMS for enhanced security
- Implement key versioning and backup
- Have disaster recovery plans for key loss

## Testing

### Unit Tests

```go
func TestEncryptDecrypt(t *testing.T) {
    key := make([]byte, 32) // test key
    enc, err := crypto.NewEncryptor(key)
    require.NoError(t, err)

    original := "test PHI data"
    encrypted, err := enc.Encrypt(original)
    require.NoError(t, err)

    decrypted, err := enc.Decrypt(encrypted)
    require.NoError(t, err)
    assert.Equal(t, original, decrypted)
}
```

### Key Validation

```go
func TestKeyValidation(t *testing.T) {
    // Invalid length
    _, err := crypto.NewEncryptor([]byte("short"))
    assert.Error(t, err)

    // Valid 32-byte key
    key := make([]byte, 32)
    _, err = crypto.NewEncryptor(key)
    assert.NoError(t, err)
}
```

## Troubleshooting

### Common Issues

1. **"PHI_ENCRYPTION_KEY must decode to exactly 32 bytes"**
   - Check hex string length (should be 64 characters)
   - Verify no extra characters or whitespace

2. **"crypto: decryption failed (authentication error)"**
   - Data may be corrupted
   - Wrong key being used
   - Check key environment variable

3. **"crypto: failed to decode encryption key"**
   - Key is not valid base64 (if using base64 format)
   - Check encoding format matches code expectations

### Recovery from Key Loss

If the encryption key is lost:
1. **Stop the application** immediately
2. **Assess data impact** - all encrypted PHI is inaccessible
3. **Restore from backup** if available (with old key)
4. **Notify affected parties** per HIPAA breach requirements
5. **Generate new key** and re-encrypt data if possible

### Monitoring

Monitor for:
- Decryption failures (potential key issues)
- Encryption operation latency
- Key rotation status
- Audit log entries for crypto operations

## Compliance Notes

- **HIPAA §164.312(a)(2)(iv)**: Encryption and decryption
- **HIPAA §164.312(e)(2)(ii)**: Encryption at rest
- Keys must be protected with same rigor as the data they protect
- Audit all cryptographic operations
- Document key lifecycle management procedures