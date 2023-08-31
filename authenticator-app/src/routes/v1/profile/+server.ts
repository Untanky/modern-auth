import type { Profile } from '$lib/profile/model';
import { profileService } from '../../../hooks.server';
import type { RequestEvent, RequestHandler } from './$types';

export const GET: RequestHandler = async ({ cookies }: RequestEvent): Promise<Response> => {
    const sub = '70827860-7316-4099-983c-4c434ca7286d';
    const profile = await profileService.find(sub);

    const headers = new Headers();
    headers.set('cache', 'no-store');
    return new Response(JSON.stringify(profile));
};

export const PUT: RequestHandler = async ({ request }: RequestEvent): Promise<Response> => {
    const sub = '70827860-7316-4099-983c-4c434ca7286d';
    const profile = await request.json() as Profile;

    await profileService.update({
        ...profile,
        sub,
    });

    const headers = new Headers();
    headers.set('cache', 'no-store');
    return new Response(null, { status: 204, headers });
};
