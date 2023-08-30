import { sendEmail } from './email-delivery';
import type {
    Email, EmailRepository, PreferencesRepository, Template,
} from './email.model';
import { renderTemplate } from './render-template';

export type EmailWithCompleteTemplate = Email & { template: Template }

export class EmailService {
    private readonly emailRepo: EmailRepository;
    private readonly preferencesRepo: PreferencesRepository;

    constructor(emailRepository: EmailRepository, preferencesRepository: PreferencesRepository) {
        this.emailRepo = emailRepository;
        this.preferencesRepo = preferencesRepository;
    }

    async send(email: EmailWithCompleteTemplate): Promise<Pick<Email, 'id'>> {
        const { body, subject } = renderTemplate(email.template);

        const { emailAddress } = await this.preferencesRepo.findFirst({ sub: email.sub });

        const resend = await sendEmail({
            to: emailAddress, body, subject,
        });
        return this.saveEmailInDB(email, resend);
    }

    private async saveEmailInDB(email: Email, resend: { id: string }): Promise<{ id: string }> {
        try {
            const createdEmail = await this.emailRepo.create({
                ...email,
                id: '',
                sentAt: new Date(),
                resendId: resend.id,
            });
            return { id: createdEmail.id };
        } catch (e) {
            return { id: '' };
        }
    }
}
