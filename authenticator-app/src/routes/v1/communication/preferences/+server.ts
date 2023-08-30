import { preferenceService } from '../../../../hooks.server';
import type { RequestHandler, RequestEvent } from './$types';

export const GET: RequestHandler = async ({ cookies }: RequestEvent): Promise<Response> => {
    const sub = 'd6952eb8-d912-4e23-ad64-d27a01a960b2';
    const preferences = await preferenceService.find(sub);

    return new Response(JSON.stringify(preferences));
};

export const PUT: RequestHandler = ({ cookies }: RequestEvent): Promise<Response> => {
    return Promise.resolve(new Response());
};
