package data

// The reported symptoms below are non specific or are not actual symptoms, so exclude them to keep the count meaningful
var ExcludeSymptoms = map[string]struct{}{
	"blood test": {},
	"laboratory test": {},
	"sars-cov-2 test positive": {},
	"covid-19": {},
	"sars-cov-2 test negative": {},
	"sars-cov-2 test": {},
	"electrocardiogram": {},
	"condition aggravated": {},
	"body temperature": {},
	"full blood count": {},
	"computerised tomogram": {},
	"metabolic function test": {},
	"chest x-ray": {},
	"electrocardiogram normal": {},
	"no adverse event": {},
	"magnetic resonance imaging": {},
	"chest x-ray normal": {},
	"heart rate": {},
	"unevaluable event": {},
	"full blood count normal": {},
	"blood pressure measurement": {},
	"intensive care": {},
	"drug ineffective": {},
	"urine analysis": {},
	"vital signs measurement": {},
	"computerised tomogram head": {},
	"echocardiogram": {},
	"exposure during pregnancy": {},
	"laboratory test normal": {},
	"resuscitation": {},
	"chest x-ray abnormal": {},
	"electrocardiogram abnormal": {},
	"endotracheal intubation": {},
	"blood test normal": {},
	"blood glucose normal": {},
	"x-ray": {},
	"influenza virus test negative": {},
	"computerised tomogram normal": {},
	"underdose": {},
	"illness": {},
	"inflammation": {},
	"troponin": {},
	"computerised tomogram head normal": {},
	"ultrasound scan": {},
	"head injury": {},
	"computerised tomogram abnormal": {},
	"metabolic function test normal": {},
	"computerised tomogram thorax": {},
	"angiogram": {},
	"troponin normal": {},
	"blood glucose": {},
	"exposure to sars-cov-2": {},
	"covid-19 pneumonia": {},
	"syringe issue": {},
	"nodule": {},
	"oxygen saturation": {},
	"echocardiogram normal": {},
	"urine analysis normal": {},
	"sars-cov-2 antibody test": {},
	"white blood cell count normal": {},
	"off label use": {},
	"platelet count normal": {},
	"ultrasound doppler": {},
	"computerised tomogram head abnormal": {},
	"magnetic resonance imaging normal": {},
	"general physical health deterioration": {},
	"pulse absent": {},
	"feeding disorder": {},
	"mechanical ventilation": {},
	"investigation": {},
	"computerised tomogram thorax abnormal": {},
	"blood potassium normal": {},
	"computerised tomogram abdomen": {},
	"blood lactic acid": {},
	"influenza": {},
	"scan with contrast": {},
	"vaccination complication": {},
	"immediate post-injection reaction": {},
	"hypersensitivity": {},
	"local reaction": {},
	"blood urea increased": {},
	"aspartate aminotransferase increased": {},
	"blood potassium decreased": {},
	"alanine aminotransferase increased": {},
	"product use issue": {},
	"pain": {},
}

/*

'Flu-like'
'Gastrointestinal'
'Psychological'
'Life threatening'
'Skin & localized to injection site'
'Muscles & bones'
'Immune system & inflammation'
'Nervous system'
'Cardiovascular'
'Eyes, mouth & ears'
'Urinary'
'Breathing'
'Balance & mobility'
'Errors by medical staff'

*/

var CategoriesMap = map[string][]string{
	"headache": {"Flu-like"},
	"pyrexia": {"Flu-like"},
	"chills": {"Flu-like"},
	"fatigue":{"Flu-like"},
	"dizziness": {"Balance & mobility"},
	"dyspnoea": {"Breathing"},
	"asthenia": {"Nervous system"},
	"malaise": {"Flu-like"},
	"cough": {"Flu-like"},
	"hyperhidrosis": {"Flu-like"},
	"oropharyngeal pain": {"Flu-like"},
	"body temperature increased": {"Flu-like"},
	"skin warm": {"Skin & localized to injection site"},
	"throat irritation": {"Flu-like"},
	"influenza like illness": {"Flu-like"},
	"lethargy": {"Flu-like"},
	"migraine": {"Nervous system"},
	"rhinorrhoea": {"Flu-like"},
	"pharyngeal swelling": {"Breathing"},
	"somnolence": {"Flu-like"},
	"nasal congestion": {"Flu-like"},
	"pneumonia": {"Flu-like", "Life threatening"},
	"discomfort": {"Nervous system"},
	"cold sweat": {"Nervous system"},
	"respiratory tract congestion": {"Flu-like"},
	"dehydration": {"Flu-like"},
	"night sweats": {"Nervous system"},
	"nasopharyngitis": {"Flu-like"},
	"head discomfort": {"Flu-like"},
	"nausea": {"Gastrointestinal"},
	"vomiting": {"Gastrointestinal"},
	"diarrhoea": {"Gastrointestinal"},
	"decreased appetite": {"Gastrointestinal"},
	"abdominal pain": {"Gastrointestinal"},
	"abdominal pain upper": {"Gastrointestinal"},
	"dysphagia": {"Gastrointestinal"},
	"abdominal discomfort": {"Gastrointestinal"},
	"retching": {"Gastrointestinal"},
	"feeling abnormal": {"Psychological"},
	"impaired work ability": {"Psychological"},
	"anxiety": {"Psychological"},
	"confusional state": {"Psychological"},
	"sleep disorder": {"Psychological"},
	"insomnia" :{"Psychological"},
	"vertigo": {"Balance & mobility"},
	"disorientation": {"Psychological"},
	"mental status changes": {"Psychological"},
	"hypersomnia": {"Psychological"},
	"sensation of foreign body": {"Psychological"},
	"nervousness": {"Psychological"},
	"loss of personal independence in daily activities": {"Psychological"},
	"death": {"Life threatening"},
	"syncope": {"Cardiovascular", "Life threatening"},
	"loss of consciousness": {"Life threatening", "Cardiovascular"},
	"cerebrovascular accident": {"Life threatening", "Cardiovascular"},
	"anaphylactic reaction": {"Life threatening", "Immune system & inflammation"},
	"atrial fibrillation": {"Life threatening", "Cardiovascular"},
	"cardiac arrest": {"Life threatening", "Cardiovascular"},
	"pulmonary embolism": {"Life threatening", "Cardiovascular"},
	"myocardial infarction":{"Life threatening", "Cardiovascular"},
	"injection site pain": {"Skin & localized to injection site"},
	"injection site erythema": {"Skin & localized to injection site"},
	"pruritus": {"Skin & localized to injection site"},
	"injection site swelling": {"Skin & localized to injection site"},
	"rash": {"Skin & localized to injection site"},
	"injection site pruritus": {"Skin & localized to injection site"},
	"erythema": {"Skin & localized to injection site"},
	"paraesthesia": {"Nervous system"},
	"injection site warmth": {"Skin & localized to injection site"},
	"urticaria": {"Skin & localized to injection site"},
	"vaccination site pain": {"Skin & localized to injection site"},
	"injection site rash": {"Skin & localized to injection site"},
	"injection site induration": {"Skin & localized to injection site"},
	"rash pruritic": {"Skin & localized to injection site"},
	"injection site reaction": {"Skin & localized to injection site"},
	"rash macular": {"Skin & localized to injection site"},
	"injection site urticaria": {"Skin & localized to injection site"},
	"vaccination site swelling": {"Skin & localized to injection site"},
	"injection site mass": {"Skin & localized to injection site"},
	"rash papular": {"Skin & localized to injection site"},
	"herpes zoster": {"Skin & localized to injection site"},
	"cellulitis": {"Skin & localized to injection site"},
	"vaccination site erythema": {"Skin & localized to injection site"},
	"contusion": {"Skin & localized to injection site"},
	"angioedema": {"Immune system & inflammation"},
	"vaccination site warmth": {"Skin & localized to injection site"},
	"myalgia": {"Muscles & bones"},
	"arthralgia": {"Muscles & bones"},
	"back pain": {"Muscles & bones"},
	"neck pain": {"Muscles & bones"},
	"mobility decreased": {"Balance & mobility"},
	"muscular weakness": {"Muscles & bones"},
	"axillary pain": {"Muscles & bones"},
	"musculoskeletal stiffness": {"Muscles & bones"},
	"gait disturbance": {"Balance & mobility"},
	"limb discomfort": {"Muscles & bones"},
	"muscle spasms": {"Nervous system"},
	"injected limb mobility decreased": {"Balance & mobility"},
	"gait inability": {"Balance & mobility"},
	"dysstasia": {"Balance & mobility"},
	"muscle tightness": {"Muscles & bones"},
	"lymphadenopathy": {"Immune system & inflammation"},
	"flushing": {"Nervous system"},
	"peripheral swelling": {"Immune system & inflammation"},
	"swelling": {"Immune system & inflammation"},
	"paraesthesia oral": {"Eyes, mouth & ears"},
	"throat tightness": {"Breathing"},
	"rash erythematous": {"Skin & localized to injection site"},
	"swelling face": {"Immune system & inflammation"},
	"hypoaesthesia oral": {"Eyes, mouth & ears"},
	"swollen tongue": {"Eyes, mouth & ears"},
	"lip swelling": {"Eyes, mouth & ears"},
	"lymph node pain": {"Immune system & inflammation"},
	"burning sensation": {"Nervous system"},
	"wheezing": {"Breathing"},
	"pallor": {"Cardiovascular"},
	"hot flush": {"Nervous system"},
	"eye swelling": {"Eyes, mouth & ears"},
	"eye pruritus": {"Eyes, mouth & ears"},
	"pain in extremity": {"Nervous system"},
	"hypoaesthesia": {"Nervous system"},
	"tremor": {"Nervous system"},
	"feeling hot": {"Nervous system"},
	"facial paralysis": {"Nervous system"},
	"feeling cold": {"Nervous system"},
	"unresponsive to stimuli": {"Nervous system"},
	"tenderness": {"Nervous system"},
	"seizure": {"Nervous system"},
	"presyncope": {"Cardiovascular"},
	"aphasia": {"Nervous system"},
	"speech disorder": {"Nervous system"},
	"dysphonia": {"Nervous system"},
	"dysarthria": {"Nervous system"},
	"balance disorder": {"Balance & mobility"},
	"oedema peripheral": {"Immune system & inflammation"},
	"heart rate increased": {"Cardiovascular"},
	"chest discomfort": {"Cardiovascular"},
	"chest pain": {"Cardiovascular"},
	"palpitations": {"Cardiovascular"},
	"tachycardia": {"Cardiovascular"},
	"blood pressure increased": {"Cardiovascular"},
	"hypertension": {"Cardiovascular"},
	"hypotension": {"Cardiovascular"},
	"hypoxia": {"Cardiovascular"},
	"peripheral coldness": {"Nervous system"},
	"vision blurred": {"Eyes, mouth & ears"},
	"visual impairment": {"Eyes, mouth & ears"},
	"photophobia": {"Eyes, mouth & ears"},
	"dysgeusia": {"Eyes, mouth & ears"},
	"ageusia": {"Eyes, mouth & ears"},
	"tinnitus": {"Eyes, mouth & ears"},
	"anosmia": {"Eyes, mouth & ears"},
	"ear pain": {"Eyes, mouth & ears"},
	"eye pain": {"Eyes, mouth & ears"},
	"dry mouth": {"Eyes, mouth & ears"},
	"pain in jaw": {"Eyes, mouth & ears"},
	"poor quality product administered": {"Errors by medical staff"},
	"product storage error": {"Errors by medical staff"},
	"product temperature excursion issue": {"Errors by medical staff"},
	"product administered to patient of inappropriate age": {"Errors by medical staff"},
	"inappropriate schedule of product administration": {"Errors by medical staff"},
	"incorrect dose administered": {"Errors by medical staff"},
	"product administration error": {"Errors by medical staff"},
	"product administered at inappropriate site": {"Errors by medical staff"},
	"incorrect route of product administration": {"Errors by medical staff"},
	"interchange of vaccine products": {"Errors by medical staff"},
	"wrong product administered": {"Errors by medical staff"},
	"fall": {"Balance & mobility"},
	"induration": {"Skin & localized to injection site"},
	"white blood cell count increased": {"Immune system & inflammation"},
	"oxygen saturation decreased": {"Cardiovascular"},
	"blood glucose increased": {"Cardiovascular"},
	"injection site bruising": {"Skin & localized to injection site"},
	"platelet count decreased": {"Cardiovascular"},
	"bell's palsy": {"Nervous system"},
	"troponin increased": {"Cardiovascular"},
	"bone pain": {"Muscles & bones"},
	"pain of skin": {"Skin & localized to injection site"},
	"taste disorder": {"Eyes, mouth & ears"},
	"oropharyngeal discomfort": {"Eyes, mouth & ears"},
	"thrombosis": {"Cardiovascular", "Life threatening"},
	"joint swelling": {"Muscles & bones"},
	"skin discolouration": {"Skin & localized to injection site"},
	"urinary tract infection": {"Urinary"},
	"injection site inflammation": {"Skin & localized to injection site"},
	"haemoglobin decreased": {"Cardiovascular"},
	"injection site nodule": {"Skin & localized to injection site"},
	"blood pressure decreased": {"Cardiovascular"},
	"epistaxis": {"Cardiovascular"},
	"ear discomfort": {"Eyes, mouth & ears"},
	"pharyngeal paraesthesia": {"Eyes, mouth & ears"},
	"thirst": {"Eyes, mouth & ears"},
	"disturbance in attention": {"Psychological"},
	"skin burning sensation": {"Skin & localized to injection site"},
	"muscle twitching": {"Nervous system"},
	"deep vein thrombosis": {"Cardiovascular", "Life threatening"},
	"blister": {"Skin & localized to injection site"},
	"neuralgia": {"Nervous system"},
	"dyspepsia": {"Gastrointestinal"},
	"restlessness": {"Psychological"},
	"joint range of motion decreased": {"Balance & mobility"},
	"lacrimation increased": {"Eyes, mouth & ears"},
	"c-reactive protein increased": {"Immune system & inflammation"},
	"sensitive skin": {"Skin & localized to injection site"},
	"fibrin d dimer increased": {"Cardiovascular"},
	"feeling of body temperature change": {"Nervous system"},
	"ocular hyperaemia": {"Eyes, mouth & ears"},
	"abdominal distension": {"Gastrointestinal"},
	"facial pain": {"Nervous system"},
	"hemiparesis": {"Nervous system"},
	"injection site cellulitis": {"Skin & localized to injection site"},
	"asthma": {"Breathing"},
	"eye irritation": {"Eyes, mouth & ears"},
	"heart rate decreased": {"Cardiovascular"},
	"sneezing": {"Flu-like"},
	"dry throat": {"Eyes, mouth & ears"},
	"breast pain": {"Nervous system"},
	"heart rate irregular": {"Cardiovascular"},
	"injection site discomfort": {"Skin & localized to injection site"},
	"hallucination": {"Psychological"},
	"sepsis": {"Life threatening", "Immune system & inflammation"},
	"respiratory arrest": {"Breathing", "Life threatening"},
	"swelling of eyelid": {"Eyes, mouth & ears"},
	"petechiae": {"Skin & localized to injection site"},
	"injection site paraesthesia": {"Skin & localized to injection site"},
	"joint stiffness": {"Balance & mobility"},
	"dyskinesia": {"Nervous system"},
	"vaccination site pruritus": {"Skin & localized to injection site"},
	"musculoskeletal chest pain": {"Muscles & bones"},
	"respiratory failure": {"Breathing", "Life threatening"},
	"acute myocardial infarction": {"Cardiovascular", "Life threatening"},
	"musculoskeletal discomfort": {"Muscles & bones"},
	"oral herpes": {"Immune system & inflammation"},
	"dyspnoea exertional": {"Breathing"},
	"hypoacusis": {"Eyes, mouth & ears"},
	"injection site hypoaesthesia": {"Skin & localized to injection site"},
	"skin reaction": {"Skin & localized to injection site"},
	"dizziness postural": {"Balance & mobility"},
	"productive cough": {"Flu-like"},
	"respiratory distress": {"Breathing"},
	"movement disorder": {"Balance & mobility"},
	"thrombocytopenia": {"Cardiovascular", "Life threatening"},
	"periorbital swelling": {"Eyes, mouth & ears"},
	"haemorrhage": {"Cardiovascular", "Life threatening"},
	"memory impairment": {"Psychological"},
	"appendicitis": {"Gastrointestinal"},
	"blood creatinine increased": {"Urinary"},
	"acute kidney injury": {"Urinary"},
}
