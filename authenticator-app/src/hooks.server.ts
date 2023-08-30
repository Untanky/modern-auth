import { db } from '$lib/server/db';
import { DrizzleEmailRepository } from '$lib/server/email.repository';
import { EmailService } from '$lib/server/email.service';
import { DrizzlePreferencesRepository } from '$lib/server/preference.repository';
import { PreferenceService } from '$lib/server/preference.service';
import { VerificationService } from '$lib/server/verification.service';


const drizzleEmailRepo = new DrizzleEmailRepository(db);
const drizzlePreferenceRepo = new DrizzlePreferencesRepository(db);
const emailService = new EmailService(drizzleEmailRepo, drizzlePreferenceRepo);
const verificationService = new VerificationService(db, emailService);
const preferenceService = new PreferenceService(drizzlePreferenceRepo, verificationService);

