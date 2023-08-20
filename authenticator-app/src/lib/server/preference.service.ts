import type { Preferences, PreferencesRepository } from './email.model';

export class PreferenceService {
    private readonly preferencesRepo: PreferencesRepository;

    constructor(preferencesRepository: PreferencesRepository) {
        this.preferencesRepo = preferencesRepository;
    }

    find(sub: string): Promise<Preferences> {
        return this.preferencesRepo.findFirst({ sub });
    }

    update(preferences: Preferences): Promise<void> {
        return this.preferencesRepo.update(preferences).then();
    }
}
