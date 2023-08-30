import type { InsertPreferences, Preferences } from '$lib/preferences/model';
import { preferenceService } from '../../../../hooks.server';
import type { RequestEvent, RequestHandler } from './$types';

export const GET: RequestHandler = async ({ cookies }: RequestEvent): Promise<Response> => {
    const sub = '70827860-7316-4099-983c-4c434ca7286d';
    const preferences = await preferenceService.find(sub);

    const headers = new Headers();
    headers.set('cache', 'no-store');
    return new Response(JSON.stringify(preferences));
};

export const PUT: RequestHandler = async ({ request }: RequestEvent): Promise<Response> => {
    const sub = '70827860-7316-4099-983c-4c434ca7286d';
    const preferences = await request.json() as InsertPreferences;

    await preferenceService.update({
        ...preferences,
        sub,
    } as Preferences);

    const headers = new Headers();
    headers.set('cache', 'no-store');
    return new Response(null, { status: 204, headers });
};
