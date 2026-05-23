package ai

// ReadDocxText uses godocx to extract all plain text from your file
var DocString string = `

LifePact Health
Complete FAQ
by Category
Official LifePact FAQ · Internal use only

How to use this document
This document is formatted for easy updating in Google Docs. Each category has a clearly labeled heading, an introductory line, and numbered Q&A entries. To add a new question, copy the format of any existing entry under the right category heading. Categories appear in the same order as the interactive online FAQ database.
01 About LifePact
Foundation questions about what LifePact is, how it works, and the role of the Wellness Advisor.
Q 1.1
What is LifePact Health?
LifePact Health is a telehealth platform providing personalized health optimization services including hormone replacement therapy (HRT), peptide therapy, GLP-1 medications, supplements, nutraceuticals, and advanced diagnostic testing. All care is managed by licensed medical providers, with dedicated Wellness Advisors guiding patients through the process.


Q 1.2
What makes LifePact different from other providers?
LifePact combines functional medicine, diagnostics, and advanced treatment options into one seamless platform. Rather than treating symptoms in isolation, LifePact focuses on optimizing how you feel, perform, and age — through comprehensive lab testing, customized treatment plans, and direct support from a personal Wellness Advisor. LifePact acts as a long-term health partner, not a quick fix.


Q 1.3
Do Wellness Advisors prescribe medication?
No. Wellness Advisors do not diagnose, prescribe, or offer medical advice. Their role is to educate patients on available services, support them throughout the program, reinforce provider guidance, answer non-medical questions, and act as a liaison between the patient and a licensed medical provider. All medical decisions are made solely by licensed healthcare providers.


Q 1.4
What is the difference between HRT and peptides?
HRT/TRT replaces declining hormones such as testosterone, estrogen, or progesterone when the body is no longer producing them at optimal levels. Peptides are short amino acid chains that signal the body to perform specific functions — such as improving fat metabolism, recovery, sleep, or cognitive performance. Both can be used alone or together depending on labs, symptoms, and goals.


Q 1.5
Do I have to stay on treatment forever?
Not necessarily. Some treatments like HRT are designed for long-term use to maintain stable hormone levels. Others, like many peptides or GLP-1s, are often used in cycles or phases. The medical provider will guide the dosing schedule based on how the body responds and what the patient is trying to achieve.


Q 1.6
What are the potential side effects?
Potential side effects vary depending on the medication but are generally mild and temporary. The provider will review all risks, benefits, and safety considerations before prescribing anything. If a patient ever experiences discomfort or an unexpected reaction, they should contact the LifePact medical team immediately.


02 Pricing & Fees
All costs, payment options, and how LifePact compares to other platforms.
Q 2.1
How much does it cost to get started?
$374 total: $199 for the comprehensive blood panel via Quest Diagnostics + $175 annual provider fee. Medication is separate and a la carte — no monthly subscriptions. Blood panel pricing may vary outside of Quest Diagnostics. Your Wellness Advisor will help explore different options if you do not have a Quest location near you.


Q 2.2
What does the $175 annual provider fee cover?
Twelve months of ongoing medical oversight — lab reviews each time you complete labs, protocol recommendations, annual Zoom consultation with a licensed provider, and ongoing Wellness Advisor support. Flat annual charge, not monthly. 


Q 2.3
What does the blood panel cost and what's included?
$199 for the comprehensive panel at Quest Diagnostics: free and total testosterone, estradiol, progesterone (females only), PSA (males only), FSH, LH, DHEA-S, IGF-1, fasted insulin, CBC, CMP, A1C, lipid panel, TSH, free T3, free T4. 


Q 2.4
Can both my husband/partner and I both sign up?
Yes — each patient pays their own $374 startup ($199 panel + $175 annual fee). Medications are a la carte for each. Many couples go through the process in parallel, and shared accountability often improves outcomes for both.


Q 2.5
Is there a fee for talking to the Wellness Advisor?
No. The Wellness Advisor is included as part of being a patient.


Q 2.6
Can I use insurance to cover any of the costs?
LifePact is cash-pay only. HSA and FSA cards may be accepted for blood work, medical provider consultations, and prescribed compounded medications — however, patients must confirm eligibility directly with their own HSA/FSA administrator before assuming coverage.


Q 2.7
What will medications typically cost?
Foundational hormones (testosterone, progesterone, DHEA) can be as little as: $200–$500 for a 4-month supply. Average protocol order: ~$1,000 for 4 months. You only pay for what you feel comfortable with and of course what you are approved for by our licensed medical team. Pricing may vary based on customized dosing, pharmacy availability, etc.


Q 2.8
How is LifePact priced compared to monthly subscription platforms?
Monthly subscription platforms like MIDI, Hone, or Wyona charge $100–$149+/month in membership fees alone — $1,200–$1,800/year before medications. LifePact's annual fee is $175/year total.


Q 2.9
What forms of payment are accepted?
Credit/debit cards, ACH payments, HSA cards, and FSA cards (pending verification with your specific administrator). Insurance is not accepted. No payment plans.


03 Process: Step by Step
What happens from the moment a patient signs up through receiving medications.
Q 3.1
What is the full step-by-step process from start to treatment?
Step 1: Complete intake form and ID verification. Step 2: Wellness Advisor consultation call. Step 3: Pay intake fee ($374 or $175 if recent labs on file). Step 4: Lab requisition processed and emailed within 1-2 days of payment, usually within a couple of hours. Step 5: Blood draw at Quest Diagnostics (or sometimes LabCorp). Step 6: Medical provider reviews labs, symptoms, history, and goals. Step 7: Lab review call with Wellness Advisor. Step 8: Pay for medication order if moving forward. Step 9: Annual Zoom consultation with a licensed provider. Step 10: Pharmacy processing and direct shipping. Step 11: Ongoing support and follow-up.


Q 3.2
How long does the full process take?
Most patients move from intake to treatment in about 7–12 business days. Labs take 3–7 business days (some markers up to 10). Provider review: 1–3 business days after labs. Medications ship in as little as 1–3 business days; some may take up to 7 business days. If no new labs are needed: as quick as 1–3 business days from intake to med team recommendations.


Q 3.3
What happens after I pay for labs?
Our operations team will process your lab order within 24-48 hours of payment. You will receive a copy of your lab requisition form via the LifePact client portal along with how to schedule your lab appointment, how to best prepare for your blood draw, and what to expect next in the process.


Q 3.4
Why do I need to upload my driver's license / ID?
ID verification is required for telehealth compliance—it confirms identity before creating a patient chart, ordering labs, or prescribing. The legal name and date of birth on the lab requisition must exactly match the current government-issued ID. Quest and LabCorp will not complete the draw if there is any discrepancy.


Q 3.5
What does the initial consultation call cover?
The initial consultation is a complimentary phone call (not Zoom) with an Onboarding Specialist or Wellness Advisor. During this call, next steps are explained, non-medical questions are answered, and goals are discussed. The Zoom call with the licensed provider comes later when moving forward with a protocol.


Q 3.6
Can I get started with just the annual fee and use my existing labs?
Yes—for patients with recent comprehensive labs (generally past 90 days), the $175 annual fee can be paid first. The Wellness Advisor reviews uploaded labs and sends them to the med team. The $199 panel is ordered when the next recheck is due. This is a common fast-track path, especially for patients transitioning from other providers.


Q 3.7
I want to proactively refill my medications before I run out—how early can I do that?
We encourage proactive refills to help you stay consistent and avoid any gaps in your protocol.
In most cases, the ideal time to request a refill is about 2–4 weeks before you run out, especially for compounded medications, since pharmacy processing and shipping may take up to 5–7 business days in some cases.
That said, refill timing can vary depending on the medication:
    • Controlled or regulated medications may only be eligible for refill within ~7–10 days of your calculated run-out date, based on your prescribed dosing.
    • Some medications also have quantity or timing limits set by the pharmacy, which can affect how early they can be reordered.
If you're unsure or trying to plan ahead for travel, schedule changes, or specific goals, your Wellness Advisor can help map out the best timing based on your protocol.


04 Labs & Bloodwork
Everything patients need to know about getting labs done, timing, and using existing results.
Q 4.1
Do I need to get labs done even if I have recent results?
Upload previous labs—they provide valuable context. If lab work is from the past 90 days and includes all required biomarkers, it can often be used. The licensed medical provider will ultimately determine whether recent labs are sufficient or if new testing is necessary.


Q 4.2
How do I get the blood draw done?
Quest Diagnostics is preferred ($199 panel vs. $399 at LabCorp). After payment, the lab requisition and booking link are sent within 1-2 days payment. Labs must be fasted, ideally done in the morning before 10am. Specific instructions on how to best prepare for your blood draw will be included in the messaging sent via the LifePact client portal.


Q 4.3
How long do lab results take?
About one week for full results (some markers up to 10 business days). The medical team review takes 1–3 business days after receipt. Total from blood draw to results call: about 1.5–2 weeks.


Q 4.4
What day of my cycle should I get labs drawn?
Ideally days 19–22 (progesterone peaks in the luteal phase) for regular cycles. For irregular cycles, go whenever—just note what cycle day it is. Days 3–4 after the period starts is another useful checkpoint. Don't wait weeks for the perfect day. Please let your Wellness Advisor know which day of your cycle you had labs completed on as that is the most important information to interpret lab results properly.


Q 4.5
How many days after my testosterone injection should I wait before labs?
About 4 days after the most recent injection for weekly dosers. Timing significantly affects readings—even 24 hours can shift numbers substantially.


Q 4.6
My previous labs showed everything was "fine" but I'm still experiencing all these symptoms—what gives?
Two root causes: (1) The wrong markers were checked—many panels miss DHEA, free testosterone, and the full thyroid panel. (2) The labs were interpreted through a disease lens, not an optimization lens. "Normal range" is very broad. LifePact looks at where you are in the range and how that aligns with your symptoms.


Q 4.7
Can biotin in my supplements skew my thyroid lab results?
Yes—biotin (vitamin B7), commonly found in multivitamins and hair supplements, can significantly interfere with thyroid readings. The classic pattern: free T3 appears very high while free T4 appears very low. Patients should stop biotin 3–5 days before labs.


Q 4.8
My lab requisition had my old last name on it and Quest wouldn't draw my blood—what do I do?
The name on the lab requisition must exactly match the patient's current legal ID. Contact the Wellness Advisor immediately—the requisition will be reprocessed as soon as possible.


Q 4.9
Should I take my thyroid medication the morning of my blood draw?
No—skip the morning dose on the day of the draw. Take it as normal the day before. Taking thyroid medication right before labs can cause a temporary spike that does not reflect the true baseline.


Q 4.10
Can I use labs from Function Health or another membership lab service?
Possibly—it depends on what markers are included. If the key biomarkers are covered, they can often substitute for the $199 panel. Upload the results and the Wellness Advisor will confirm what is covered and what needs to be added.


Q 4.11
My existing labs are missing key markers—what can be done?
LifePact can build a custom partial panel to fill in specific missing markers rather than re-ordering the entire $199 comprehensive panel, in some cases. The Wellness Advisor will price a targeted add-on and communicate the cost before ordering. 


05 Hormones & HRT
Questions about hormone replacement therapy including testosterone, estrogen, progesterone, DHEA, and thyroid.
Q 5.1
Do you offer pellet therapy?
No. Pellet release is unpredictable and cannot be adjusted once inserted. LifePact uses subcutaneous injections and topical creams for HRT in most cases for precise, adjustable dosing.


Q 5.2
What forms of testosterone does LifePact offer for women?
Primarily subcutaneous injection (small insulin syringe, once weekly). Topical cream is also available but subQ tends to be preferred—more bioavailable, more precise, and less likely to convert to DHT. The pros and cons of different forms of administration can always be discussed with your Wellness Advisor or provider.


Q 5.3
What is estrogen's role in histamine and itchy skin / ears?
When estrogen declines, it can trigger histamine release. Histamine breakdown is what causes itchiness — particularly in ears and scalp — that worsens in the luteal phase (week before the period). This is not allergies; it's an estrogen-histamine interaction.


Q 5.4
Why do I wake up at 3am—is that hormonal?
Yes, it usually is. When estrogen and progesterone are not balanced, the body's internal thermostat gets dysregulated. This can trigger cortisol or adrenaline surges during sleep—which jolts you awake at 3am. Progesterone support has a calming, thermogenic effect that helps stabilize sleep architecture.


Q 5.5
I had a hysterectomy but still have my ovaries—does that change my hormonal situation?
Yes—significantly. Keeping the ovaries means they still produce estrogen, progesterone, and testosterone. This is very different from a full hysterectomy where ovaries are removed. Without a uterus, there are no periods as a tracking signal, but hormonal decline still happens gradually as the ovaries age.


Q 5.6
I have a family history of breast cancer—can I still do hormone therapy?
A family history of breast cancer requires careful, individualized evaluation. In most cases, hormone therapies such as estrogen, DHEA, or testosterone are not typically recommended without a very clear clinical need and thorough review. If there is ever a situation where hormone therapy may be considered appropriate, our medical team would also require clearance from your oncologist before moving forward.
That said, this does not mean you’re out of options.
There are several non-hormonal and lower-risk approaches that can still effectively support symptoms, optimize health, and improve quality of life. Your protocol would be built around what is both safe and appropriate for your specific history.
Our medical team takes your full personal and family history into account to determine the best path forward.


Q 5.7
I developed acne after starting testosterone—is that from the hormone or from menopause?
Both are possible.
Testosterone can contribute to acne in some women, especially depending on the dose, delivery method (like creams), and individual sensitivity. At the same time, hormonal shifts during menopause can also increase breakouts due to changes in estrogen, progesterone, and skin oil production.
The important thing to know is that this is something we can absolutely adjust and manage.
We take a conservative, individualized approach from the start, typically using lower, patient-specific dosing to minimize the risk of side effects like acne. If breakouts do occur, it does not automatically mean testosterone is the wrong fit—it often just means we need to fine-tune the dose, frequency, or delivery method.
Also, not all acne is purely testosterone or androgen-driven. Factors like skin sensitivity, inflammation, stress, and overall hormonal balance can all play a role.
If you notice changes, we’ll work with you to adjust your protocol so you can continue to see benefits without unwanted side effects.


Q 5.8
My doctor says my thyroid is "normal" but I still feel terrible—can you help?
Yes—extremely common. LifePact evaluates free T3 and free T4 (not just TSH) through an optimization lens. "In range" doesn't mean optimal.


Q 5.9
Do you offer thyroid support / medication?
Yes. Compounded desiccated thyroid (equivalent to Armour Thyroid) in capsule form and synthetic options. The compounded version is often more accessible and affordable for most.


Q 5.10
I have Graves' disease—can LifePact help me switch to Armour Thyroid?
Yes. Patients with Graves' disease on levothyroxine (T4 only) often find that switching to desiccated thyroid (compounded equivalent of Armour) produces better results because it contains both T3 and T4. The transition requires careful lab monitoring.


Q 5.11
I ran out of my thyroid medication from my previous provider—how fast can LifePact help?
Thyroid medication is treated as a priority, and we work quickly to avoid any disruption in your care. In many cases, we can process and submit your order the same day, as long as your recent labs and onboarding forms are completed and up to date.
From there, timing depends on the pharmacy, but we can often expedite processing and shipping when needed. In some cases, this may involve an additional fee, depending on how quickly the medication is required.
If you’re running low or already out, the best step is to reach out as soon as possible so we can move quickly and coordinate the fastest option available.


06 Perimenopause Symptoms
Understanding the hormonal mechanisms behind common perimenopause symptoms.
Q 6.1
Everyone is telling me I'm too young for perimenopause—but something is clearly wrong.
You know your body. Perimenopause can start in the mid-30s, and chronic stress can trigger hormonal changes regardless of age. FSH, LH, estradiol, and progesterone levels don't care what the calendar says. LifePact doesn't need to put a label on it—the goal is to identify what's happening and address it.


Q 6.2
What is luteal phase instability and how does it cause my symptoms?
The luteal phase is days 14–28 of the cycle—when progesterone should rise and peak around days 19–22. When progesterone doesn't rise adequately, it destabilizes the hormonal cascade: sleep gets disrupted, temperature regulation goes haywire (night sweats, hot flashes), histamine rises (itching, anxiety), mood becomes volatile, and libido drops.


Q 6.3
I notice my worst symptoms happen the week before my period—what does that mean?
This is one of the clearest diagnostic signals LifePact receives. Symptoms peaking in the week before the period and resolving once bleeding starts indicate that progesterone is insufficient or dropping too quickly during that phase. The solution is often progesterone support timed to days 14–28 of the cycle.


Q 6.4
How does night sweating / drenched in sweat at night relate to hormones?
Night sweats in the luteal phase are driven by the interplay between progesterone (which has a thermogenic effect) and cortisol surges that occur when hormones are imbalanced. Progesterone support in the luteal phase is typically the most effective intervention.


07 GLP-1s & Weight Loss
Everything about semaglutide, tirzepatide, dosing philosophy, and what LifePact does and does not offer.
Q 7.1
What's the difference between semaglutide and tirzepatide?
Semaglutide (Ozempic/Wegovy): first-generation GLP-1. Tirzepatide (Mounjaro/Zepbound): second-generation, also targeting GIP receptors—often produces better results with fewer side effects at comparable doses.


Q 7.2
What is microdosing / custom lower dosing for GLP-1s?
Custom patient-specific dosing starting well below the standard titration protocol (1–2mg/week vs. the standard 2.5mg start). Goal: extract the benefits—reduced food noise, anti-inflammatory effects, insulin sensitivity—without harsh side effects.


Q 7.3
I had terrible side effects on Wegovy / Ozempic—can you still help?
Yes. Side effects at standard doses can usually be avoided or minimized with custom lower dosing or by switching to tirzepatide. Many patients who stopped due to tolerability find that a restart at 1mg is a completely different experience. This doesn’t guarantee these options may be the best for you, but our medical team will evaluate what may be best for you, and we will help you navigate through that.


Q 7.4
I lost weight on tirzepatide at 2.5mg but it was too aggressive—can I restart at a lower dose?
Yes—this is exactly the scenario where custom lower dosing shines. Restarting at 1mg or 0.5mg/week typically produces a tolerable and sustainable experience while still delivering the food noise and anti-inflammatory benefits.


Q 7.5
I reached my goal weight on tirzepatide—what now?
Options: (1) Maintain at a low dose to preserve metabolic and anti-inflammatory benefits. (2) Taper down or cycle on/off. (3) Discontinue if lifestyle habits are well established. LifePact can also re-evaluate the full hormonal picture—patients who have lost significant weight often find new optimization opportunities.


Q 7.6
The main benefit I notice from tirzepatide is that the food noise is gone—is that normal?
Yes—food noise reduction is one of the most consistently reported effects at lower doses. GLP-1 medications work on receptors in the brain as well as the gut, which is why this effect often feels remarkably complete.


Q 7.7
Does LifePact offer retatrutide (third-generation GLP-1)?
No—retatrutide is still pending FDA approval and cannot be legitimately compounded through licensed 503A or 503B pharmacies. LifePact only offers medications that can be legally and properly compounded


Q 7.8
I switched from tirzepatide to retatrutide and don't feel as good—is that common?
Yes, it can happen. Patients who respond well to tirzepatide at low doses and switch to retatrutide often feel less effective overall. Part of this may be sourcing quality. It’s completely normal for different patients to respond to different medications or variations of that medication somewhat differently, but our team will help you navigate through all of this.


Q 7.9
I had an allergic reaction (angioedema) to a GLP-1—can I still work with LifePact?
Yes—LifePact can still help with hormonal optimization and other peptides. GLP-1 medications must be avoided and clearly documented in the health history. Metabolic support can be addressed through other pathways.


08 Growth Hormone Peptides
Tesamorelin, CJC-1295, Sermorelin, Hexarelin, and how growth hormone support is chosen.
Q 8.1
What's the difference between Tesamorelin and Sermorelin?
Both stimulate the pituitary to produce growth hormones naturally. Tesamorelin is more targeted and potent—specifically FDA-approved for visceral fat reduction. Sermorelin is milder, more affordable, better for general anti-aging, recovery, sleep quality, and longevity. IGF-1 levels on the blood panel are the key marker for deciding which and how much.


Q 8.2
What is CJC-1295 and how does it compare to Hexarelin?
CJC-1295 is a growth hormone-releasing hormone analogue—it produces a more sustained GH pulse and is generally better tolerated over longer protocols than Hexarelin. CJC-1295 is often paired with Ipamorelin for synergistic effect and is generally the preferred upgrade from Hexarelin.


Q 8.3
I was on Tesamorelin before and had sleep disturbances—is that normal?
Yes, it can happen. Sleep disturbances are a known side effect of Tesamorelin, especially at higher doses or when taken too close to bedtime. Best practice: inject in the evening but a few hours before bed. If sleep disruption persists, Sermorelin is a milder alternative.


Q 8.4
I was on Tesamorelin daily—is that the correct dosing frequency?
Tesamorelin is typically dosed once daily. It has a short half-life but stimulates a natural growth hormone pulse, which is why consistent daily use is standard in most clinical protocols.Some providers may choose to cycle it (such as 5 days on, 2 days off), but this is based on individual preference or patient response.


09 Peptides & Supplements
BPC-157, NAD+, NMN, 5-amino-1MQ, hair loss peptides, and the difference between research-grade and pharmaceutical-grade.
Q 9.1
Do you offer BPC-157 / PDA (Pentadecapeptide Arginate)?
Yes. Injectable BPC-157 (now often called Pentadeca Peptide or PDA in compounded form) is available through pharmacies shipping to applicable states. Used primarily for tissue repair, tendon/ligament healing, pain reduction, and gut health. Depending on the pharmacy and prescribed protocol, it generally ranges in price between $750-$1500 for a 2-4 Month protocol.


Q 9.2
How does LifePact prioritize what to recommend—do I need a lot of things?
No, not necessarily! And this is something we take seriously.
You can expect an honest, straightforward assessment of what we’re actually seeing in your labs, symptoms, and goals. From there, we help you prioritize a clear game plan based on what is most likely to move the needle first.

Even if the most “ideal” or comprehensive protocol isn’t within your budget or comfort level, we will always walk you through what the minimum effective approach could look like along with realistic expectations for what that level of intervention can achieve.
In many cases, the most impactful interventions are not the most expensive. We regularly help patients build protocols that align with a wide range of budgets while still driving meaningful results.
Certain peptides and more advanced therapies are layered in when appropriate, but only when they truly make sense for your situation. In fact, it’s common for our team to recommend starting simpler before adding complexity.
We pride ourselves on building trust through transparency and education, so you feel confident in your plan, understand your options, and know exactly why each recommendation is being made.


Q 9.3
Is NAD+ worth it—should I be taking it?
NAD+ declines with age and stress. LifePact's recommended approach: start with NMN (a precursor converted to NAD+ inside the cell), which is oral, more affordable, and better at sustaining NAD+ levels long-term. Use injectable NAD+ as an occasional accelerator.
Note: We have a great comparison document of NAD+ and NMN that is already uploaded in our CRM, so it can be easily shared with patients, as well.


Q 9.4
What's the difference between NAD+ and NMN?
NMN (nicotinamide mononucleotide) is a precursor to NAD+. Taken orally, it gets converted into NAD+ inside cells — more effective at raising intracellular NAD+ levels long-term. NAD+ injectable gives a noticeable short-term boost but is more surface-level and expensive. Strategy: NMN daily as the foundation, injectable NAD+ occasionally as a booster.
Note: We have a great comparison document of NAD+ and NMN that is already uploaded in our CRM, so it can be easily shared with patients, as well.


Q 9.5
What is 5-amino-1MQ and what is it used for?
5-amino-1MQ is an NNMT inhibitor—it blocks an enzyme linked to fat cell expansion and metabolic dysfunction. It promotes fat loss, improves metabolic function, and supports body composition. It comes in 50mg capsules taken 2–3 times daily. It can be used alongside GLP-1s or as an alternative for patients who cannot use GLP-1s.


Q 9.6
What treatments does LifePact offer for hair loss?
LifePact approaches hair loss from both a scalp and internal perspective.
Topically, we may recommend custom blends with ingredients like minoxidil and GHK-Cu to support follicle health and stimulate growth (typically ~$100–$150+). That said, the most important step is understanding the root cause. Hair loss is often influenced by factors like DHT, thyroid function, iron/ferritin, cortisol, and key nutrient levels.
For that reason, we may recommend reviewing or expanding labs to ensure we’re addressing the underlying drivers—not just the symptom. Our goal is to build a targeted plan that supports both regrowth and long-term hair health, not just a temporary fix.


Q 9.7
Should I cycle on and off peptides, or can I take them continuously?
Depends on the peptide. Growth hormone peptides (Tesamorelin, Sermorelin, CJC-1295) are typically cycled to prevent receptor desensitization. Hormones like testosterone and progesterone are generally continuous. GLP-1s can be used continuously or cycled strategically. It depends on the peptide, protocol, and overall context. 


10 Research vs. Medical Peptides
Why LifePact only uses licensed compounding pharmacies and what that means for safety and quality.
Q 10.1
What's the difference between research-grade peptides and the peptides LifePact prescribes?
Research-grade peptides are not manufactured for human use. They are produced in facilities not required to follow medical-grade safety, sterility, purity, or potency standards. LifePact exclusively works with licensed 503A and 503B compounding pharmacies that verify purity and potency, provide lot tracking, and comply with FDA, USP, and state pharmacy oversight.


Q 10.2
Why are research peptides often cheaper?
Research peptides are sold "for research use only," which allows them to bypass the medical, legal, and safety standards required for human use. They do not undergo sterility testing, potency verification, or stability testing. LifePact maintains competitive pricing while still prioritizing quality and medical oversight.


Q 10.3
Are research-grade peptides unsafe?
Research-grade peptides may contain unknown ingredients, incorrect dosages, impurities, or contaminants. Because they are not made for human use, there is no guarantee of sterility, stability, or accuracy. Compounded peptides from licensed pharmacies go through sterility and potency testing, follow USP 797 guidelines, and are dispensed by prescription only.


Q 10.4
Is it legal to take research-grade peptides?
Research-grade peptides are not approved for human consumption and are typically labeled "not for human use" or "for research purposes only." Some companies use "physician use only" labeling which may appear legitimate but does NOT mean the product meets medical pharmacy standards—it is a marketing tactic. LifePact prescriptions include full traceability: patient name, prescribing provider, pharmacy lot number, exact medication strength, dosing instructions, and tamper-evident seal.


Q 10.5
Can LifePact prescribe something if I've already been using research peptides?
Yes—patients can still join LifePact. The medical provider will need to complete a full clinical review before determining what is safe or appropriate going forward.


11 Men's Health & TRT
Testosterone replacement therapy, optimal ranges, HCG, Clomiphene, and what LifePact checks that most providers don't.
Q 11.1
I'm already on TRT from my regular doctor — why would I switch to LifePact?
PCPs check narrow markers (total testosterone, hematocrit) and focus on keeping levels "in range." LifePact checks DHEA, IGF-1, full thyroid panel, fasting insulin, and more. Plus access to peptides, HCG, Clomiphene, and other tools PCPs don't offer. Optimal male total testosterone is 800–1,200 — most clinic-managed TRT patients are in range but below optimal.


Q 11.2
My testosterone is 411—which is "normal"—but is that actually optimal?
Not by LifePact's standards. The conventional "normal" range for male patients (260–916) is based on population averages that include men with suboptimal health. LifePact targets approximately 800–1200 for total testosterone and the upper third of the reference range for free testosterone—where most men report significantly better energy, body composition, recovery, and mood. We are not simply chasing numbers, but we also take into account symptoms and other details about the individual patient’s health history, etc.


Q 11.3
What is Enclomiphene and how is it different from TRT for a younger man?
Enlomiphene signals the pituitary to produce more LH and FSH, which stimulates the testes to produce more natural testosterone. Unlike TRT, it does not suppress natural production, does not create dependency, and preserves fertility. Often preferred for younger men who want to maintain natural production before committing to TRT.


Q 11.4
If I start TRT, do I have to stay on it forever?
For men, the answer is often yes long-term—but the starting point, dosing, and duration all matter.
When starting TRT, you should be prepared for it to be an ongoing commitment. Over time, your body’s natural production can downregulate, but this does not happen immediately or overnight, and the degree of suppression can vary depending on the individual and protocol.
For women, the context is different. Many women use testosterone for a period of time to support symptoms, and depending on their situation, may not need to stay on it indefinitely. Hormones, lifestyle, and goals can shift over time.
There are several scenarios that can influence long-term use, which is why we focus on providing personalized guidance. We’ll always walk you through what to expect based on your specific situation so you can make an informed decision with the right expectations.


Q 11.5
My labs show high cholesterol and elevated hemoglobin/hematocrit on TRT—what does that mean?
In patients on testosterone, these three findings together are almost always caused by testosterone running too high. Testosterone increases red blood cell production, raising hematocrit and hemoglobin, which can affect cholesterol. All three typically resolve by adjusting the testosterone dose downward. Specific scenarios will be fully evaluated during the medical team review for the patient, though.


Q 11.6
Do you offer HCG for men on TRT?
Yes. HCG can be used alongside testosterone to help support testicular function and maintain aspects of natural hormone signaling. It’s often included to help preserve fertility and reduce the likelihood of testicular atrophy that can occur with TRT alone.
That said, it’s not required for every patient. Whether HCG is included depends on your goals, baseline labs, and overall protocol.


Q 11.7
Do you offer Cialis for men?
Yes. Tadalafil (Cialis) is available through LifePact's licensed compounding pharmacies. It's commonly used by men on TRT for both sexual function and cardiovascular benefits.


Q 11.8
I've been self-managing my TRT without a provider—can LifePact take over my care?
Yes. LifePact can transition patients who have been self-managing TRT to a fully monitored, properly prescribed protocol. This involves a current comprehensive blood panel to establish a baseline, reviewing the existing protocol, and moving to pharmacy-sourced compounded testosterone with a proper prescription.


12 Patient Eligibility
Who can use LifePact, special situations, and additional diagnostic options.
Q 12.1
Who is eligible to use LifePact Health?
LifePact is available to U.S. residents nationwide with very few prescribing and telehealth limitations. Availability of certain medications may vary by state due to pharmacy regulations.


Q 12.2
What if I'm already receiving treatment elsewhere?
Patients can absolutely continue their care with LifePact. The medical provider will review current medications, past protocols, and treatment history, then determine whether it is appropriate to continue, adjust, or optimize the existing regimen. Do not stop any current treatment unless a licensed medical provider instructs you to do so.


Q 12.3
Can I use LifePact Health if I'm pregnant or trying to conceive?
Patients can still become a LifePact patient, but certain treatments will not be appropriate. Our medical providers will only recommend options that are safe and clinically appropriate.


Q 12.4
What additional diagnostic tests does LifePact offer beyond the standard blood panel?
LifePact offers access to a wide range of advanced diagnostic testing when additional insight is needed. This may include options such as DUTCH hormone testing, salivary cortisol panels, micronutrient analysis, and comprehensive gut health testing, among others.
Not every patient needs advanced testing, but when appropriate, these tools can help us gather deeper data and build a more targeted plan. Costs vary depending on the test and lab used.


Q 12.5
I was over-prescribed DHEA or another hormone without proper monitoring from another provider and had adverse effects—what should I do now?
This can happen, especially when protocols aren’t properly monitored over time. If you’re coming from another provider, we’ll start by gathering key details around your previous protocol—including dosing, timing, duration, prior labs, and any side effects you experienced.
From there, we’ll determine whether updated labs are needed to properly assess your current levels, and make recommendations on what should be adjusted, continued, or discontinued altogether.
DHEA and pregnenolone in particular can accumulate and lead to androgenic side effects if not monitored appropriately, which is why ongoing lab tracking and proper dosing adjustments are essential.
In many cases, we’re able to refine or simplify protocols—sometimes with just small changes—to create a much better overall response and experience moving forward.


13 Couples Signing Up Together
How the process works when both partners want to optimize.
Q 13.1
My husband/wife and I both want to sign up—how does that work?
Both partners sign up as separate patients, each with their own $374 startup (or just the annual fee if recent labs exist). The Wellness Advisor can set both up on the same call. Many couples find that going through the process together improves accountability and lifestyle alignment.


Q 13.2
My husband got on TRT and is feeling amazing—I want to feel that way too.
This is one of the most relatable entry points. When one partner gets optimized and the other doesn't, the gap becomes very noticeable—in energy, libido, mood, and motivation. LifePact helps women get the same level of optimization their male partners often access more readily.


14 Medications & Pharmacy
How medications are delivered, what forms are available, and what 503A and 503B pharmacies mean.
Q 14.1
Can I pick up my prescription at a regular pharmacy like CVS or Walgreens?
No. LifePact prescriptions are filled through licensed compounding pharmacy partners. These pharmacies provide customized doses and formulations that are not available at retail pharmacies like CVS or Walgreens. We’ve also thoroughly vetted our pharmacy partners, so you don’t have to navigate that process on your own. Once you’re ready to move forward, we handle everything—from invoicing to sending your prescription in for fulfillment directly through our network.


Q 14.2
What is the difference between SubQ and IM injections?
SubQ (subcutaneous): injected just under the skin (typically abdomen or thigh). Absorption is steady and gradual, and this is the most common method for many hormones and peptides.
IM (intramuscular): injected deeper into the muscle (such as the glute or deltoid). Absorption is typically faster.
Both methods are safe and can be self-administered with proper guidance. Some medications are better suited for SubQ vs IM, and we’ll always walk you through what’s most appropriate and what to expect based on your specific protocol.


Q 14.3
What are the different ways medications can be administered?
Not all treatments require injections. Options may include: injections (SubQ or IM), oral tablets or capsules, topical creams or gels, nasal sprays, and dissolvable troches placed between the gumline and cheek or under the tongue.


Q 14.4
What is a troche and how does it work?
A troche is a small, dissolvable tablet that’s placed between the gumline and cheek or under the tongue. It dissolves slowly and allows the medication to be absorbed directly into the bloodstream, bypassing the digestive system. This can improve absorption for certain medications. Troches are typically square-shaped and can be easily cut into halves or quarters to support precise dosing when needed.


Q 14.5
What is the difference between a 503A and 503B compounding pharmacy?
503A pharmacies compound medications for individual patients based on a specific prescription—the most common type LifePact uses. 503B pharmacies (outsourcing facilities) can produce larger batches under FDA oversight and are held to equally as strict manufacturing standards. Both follow strict sterility, potency, and safety protocols.


Q 14.6
I want to consolidate all my medications to one provider—is that possible?
In most cases, yes. While we don’t prescribe every medication, if it falls within health optimization, hormones, peptides, or supplements, we can typically help consolidate your care through LifePact. Having everything managed in one place allows the medical team to see how your full protocol works together. When treatments are split across multiple providers, important details can be missed, which often leads to less optimal outcomes.
Our goal is to simplify your care, improve coordination, and ensure everything is aligned toward the same result.


15 Competitors & Comparison
How LifePact compares to Transcend, Joy and Blokes, and other platforms.
Q 15.1
How is LifePact different from other telehealth HRT/wellness platforms?
Most platforms are either monthly subscription services with formulaic protocols or self-serve order portals with minimal guidance. LifePact is guided, relationship-based, with labs interpreted through an optimization lens, protocols customized and adjusted over time, and a Wellness Advisor as the ongoing main point of contact.


Q 15.2
I was with another provider and wasn't happy—what's different?
Communication and responsiveness are LifePact's core differentiators. The Wellness Advisor is the single consistent point of contact—not a rotating team. If hormones felt good initially and then started slipping back, that's a monitoring failure, not a treatment failure. LifePact's feedback loop means how you're feeling drives protocol adjustments.



16 Mindset & Lifestyle
The philosophy behind LifePact's approach and how to frame optimization vs. treatment.
Q 16.1
I've been doing everything right—tracking macros, lifting, calorie deficit—but I'm still not losing weight.
This is the #1 thing LifePact hears and the clearest signal that something physiological is the bottleneck—most often hormonal imbalances, thyroid dysfunction, insulin resistance, or cortisol elevation from chronic stress. This is typically not a willpower problem. Labs identify the specific deficiency.


Q 16.2
My doctor told me "this is just how it is at your age"—is that true?
No. The fact that hormones naturally decline with age doesn't mean suffering through the consequences is inevitable. LifePact doesn't accept "you're getting older" as a treatment plan.


Q 16.3
Hormones are all over social media—how do I know what's real vs. hype?
Social media has raised awareness, but it’s also created a lot of noise.

At LifePact, we don’t base decisions on trends—we start with objective data. Your labs show what’s actually happening in your body, and we pair that with your symptoms, history, goals, and what has or hasn’t worked for you. This allows us to cut through the hype and build a plan that’s specific, grounded, and actually relevant to you.


Q 16.4
I don't want to be put on birth control—I want real answers about my hormones.
This is exactly the gap LifePact fills. Birth control is frequently offered as a default "solution" for hormone-related symptoms when it is not actually a solution—it masks symptoms by suppressing hormonal cycles rather than addressing the underlying imbalance.


Q 16.5
My stress caused me to gain 25 pounds even though I was still working out—is that hormonal?
It can be, but there are often multiple factors at play. Chronic stress can drive elevated cortisol, which may impact fat storage, muscle retention, blood sugar regulation, and even downstream hormones like thyroid and sex hormones. From there, you can see a cascade effect—changes in energy, recovery, appetite, sleep, and overall metabolism.
That said, stress is just one piece of the puzzle. There could be several contributing factors happening at the same time, which is exactly why establishing a proper baseline is so important.
A comprehensive blood panel helps us see how these patterns are showing up in your body specifically, so we’re not guessing—we’re identifying what’s actually driving the changes and building a plan from there.


Q 16.6
I don't care about looking 20—I just want to feel 20 and not look 50. Is that a reasonable goal?
Completely reasonable—and the exact philosophy LifePact is built around. Performance, energy, recovery, cognitive sharpness, and quality of life are the primary targets. Looking better is often a byproduct, not the goal.


Q 16.7
Women's health and hormones feel like a topic that medicine is decades behind on—why is that?
Because historically, clinical research and medical education were built around male physiology. The hormonal complexity of the female cycle, perimenopause, and menopause was largely undertreated and underresearched. Most PCPs are trained in disease management, not wellness optimization. This is the exact gap that LifePact and the broader functional medicine movement are working to address.


18 Glossary & Definitions
Key terms every advisor and patient should understand.
Q 18.1
What is a Wellness Advisor?
A non-clinical support role that helps guide patients through onboarding, treatment plans, and lifestyle integration. Wellness Advisors do not diagnose, prescribe, or offer medical advice—they serve as the primary point of contact between the patient and the licensed medical team.


Q 18.2
What is a Medical Provider at LifePact?
A licensed medical provider (physician, NP, or PA) who reviews labs, diagnoses, and prescribes treatment. All medical decisions at LifePact are made solely by licensed healthcare providers.


Q 18.3
What are peptides?
Chains of amino acids that signal the body to perform specific functions such as growth hormone release, anti-inflammatory effects, fat loss, or recovery support. Peptides are distinct from hormones—they signal the body to do something, rather than directly replacing what the body produces.


Q 18.4
What is a lab requisition?
A form required to get blood drawn at an approved diagnostic lab. LifePact processes and emails the lab requisition after the patient pays the startup invoice. The patient's name and date of birth on the requisition must exactly match their government-issued ID.


Q 18.5
What is a compounding pharmacy?
A pharmacy that customizes medication dosages and delivery forms per patient needs. LifePact exclusively uses licensed 503A and 503B compounding pharmacies. These medications are not available at retail pharmacies like CVS or Walgreens.


Q 18.6
What are GLP-1 medications?
A class of medications that mimic natural gut hormones to support fat loss, appetite control, and insulin regulation. Examples include semaglutide (Ozempic/Wegovy) and tirzepatide (Mounjaro/Zepbound). LifePact offers compounded versions at competitive pricing through licensed pharmacy partners.


Q 18.7
What are troches?
Dissolvable tablets placed under the tongue or between the gumline and cheek for absorption directly into the bloodstream. They bypass the digestive system, which can improve bioavailability. Some hormone formulations are available as troches as an alternative to injections or creams.
`
