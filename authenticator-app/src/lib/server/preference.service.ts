import type { Preferences, PreferencesRepository } from './email.model';
import type { VerificationService } from './verification.service';

export class PreferenceService {
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

        await this.#preferencesRepo.update(preferences);

        if (oldPreferences.emailAddress !== preferences.emailAddress) {
            console.log('Send verification email');
            await this.#verificationService.startVerification(preferences.sub);
        }
    }
}
