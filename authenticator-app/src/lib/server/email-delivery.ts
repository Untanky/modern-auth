import { Resend } from 'resend';
import { env } from '$env/dynamic/private';

const resend = new Resend(env.RESEND_API_KEY);

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
    from: 'lukaskingsmail@gmail.com',
    to: params.to,
    html: params.body,
    subject: params.subject,
  });
};
