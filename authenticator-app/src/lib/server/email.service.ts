import { sendEmail } from './email-delivery';
import type { Email, EmailRepository, PreferencesRepository } from './email.model';
import { renderTemplate } from './render-template';

export class EmailService {
  private readonly emailRepo: EmailRepository;
  private readonly preferencesRepo: PreferencesRepository;
  
  constructor(emailRepository: EmailRepository, preferencesRepository: PreferencesRepository) {
    this.emailRepo = emailRepository;
    this. preferencesRepo = preferencesRepository;
  }
  
  async send(email: Email): Promise<Pick<Email, 'id'>> {
    const preferences = await this.preferencesRepo.findFirst({ sub: email.sub });
    const { body, subject } = renderTemplate(email.template);
    
    const { id: resendId } = await sendEmail({
      to: preferences.email,
      body,
      subject: subject,
    });

    const createdEmail = await this.emailRepo.create({
      ...email,
      id: '',
      sentAt: new Date(),
      resendId,
    });

    return { id: createdEmail.id };
  }
}