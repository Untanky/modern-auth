import type { Preferences } from '$lib/preferences/model';
import type { PreferencesRepository } from './email.model';
import type { VerificationService } from './verification.service';

export class PreferencesService {
    readonly #preferencesRepo: PreferencesRepository;
    readonly #verificationService: VerificationService;

    constructor(
        preferencesRepository: PreferencesRepository,
        verificationService: VerificationService,
    ) {
        this.#preferencesRepo = preferencesRepository;
        this.#verificationService = verificationService;
    }

    find(sub: string): Promise<Preferences> {
        return this.#preferencesRepo.findFirst({ sub });
    }

    async update(preferences: Preferences): Promise<void> {
        const oldPreferences = await this.find(preferences.sub);

        const emailChanged = oldPreferences.emailAddress !== preferences.emailAddress;
        if (emailChanged) {
            preferences.verified = false;
            preferences.verifiedAt = undefined;
        }
        await this.#preferencesRepo.update(preferences);

        if (emailChanged) {
            await this.#verificationService.startVerification(preferences.sub);
        }
    }
}
