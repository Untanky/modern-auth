import { Resend } from 'resend';

const resend = new Resend('abc');

export type EmailDeliverable = {
  to: string;
  body: string;
  subject: string;
};

type EmailResponse = {
  id: string;
}

export const sendEmail = (params: EmailDeliverable): Promise<EmailResponse> => {
    return resend.sendEmail({
        from: 'example@lukasgrimm.me',
        to: params.to,
        html: params.body,
        subject: params.subject,
    });
};
