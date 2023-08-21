import { db } from '$lib/server/db';
import type { PreferencesRepository } from '$lib/server/email.model';
import { DrizzleEmailRepository } from '$lib/server/email.repository';
import { EmailService } from '$lib/server/email.service';


const drizzleEmailRepo = new DrizzleEmailRepository(db);
const emailService = new EmailService(drizzleEmailRepo, null as unknown as PreferencesRepository);

await emailService.send({
    deliveryMethod: 'resend',
    id: '',
    resendId: '',
    sub: '3b618858-4957-4d9e-9e54-3fe15479574f',
    template: {
        type: 'verification',
        props: {
            doiLink: 'http://localhost:3000/v1/verification/abc',
            userName: 'Untanky',
            lastName: 'Grimm',
        },
    },
});
