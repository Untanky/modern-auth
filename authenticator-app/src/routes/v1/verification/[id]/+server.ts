import { verificationService } from '../../../../hooks.server';
import type { RequestEvent, RequestHandler } from './$types';

export const GET: RequestHandler = async ({ params, url }: RequestEvent): Promise<Response> => {
    const { id } = params;
    const code = url.searchParams.get('code');
    if (!code) {
        return new Response(JSON.stringify({
            error: 'BAD_REQUEST',
            description: 'missing parameter',
        }), { status: 400 });
    }

    await verificationService.finishVerification(id, code);

    const headers = new Headers();
    // TODO: redirect to a better page
    headers.set('cache', 'no-store');
    return new Response('ok', { headers });
};
