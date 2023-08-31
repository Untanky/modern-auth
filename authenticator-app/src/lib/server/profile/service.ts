import type { Profile } from '$lib/profile/model';
import type { ProfileRepository } from '../email.model';

export class ProfileService {
    readonly #profileRepo: ProfileRepository;

    constructor(profileRepository: ProfileRepository) {
        this.#profileRepo = profileRepository;
    }

    find(sub: string): Promise<Profile> {
        return this.#profileRepo.findFirst({ sub });
    }

    update(profile: Profile): Promise<Profile> {
        return this.#profileRepo.update(profile);
    }
}
