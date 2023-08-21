import { db } from '$lib/server/db';
import type { PreferencesRepository, Template } from '$lib/server/email.model';
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
        sub: '3b618858-4957-4d9e-9e54-3fe15479574f',
    } as Template,
});
