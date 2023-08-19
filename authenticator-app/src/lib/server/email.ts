import { Resend } from 'resend';
import { renderTemplate, type Template } from './templates';
import { env } from '$env/dynamic/private';

const resend = new Resend(env.RESEND_API_KEY);

export interface EmailParams<TemplateType extends Template> {
  /*
    The email address to send the email to
   */
  to: string;
  /*
    The email template to use for sending the email
   */
  template: TemplateType['template'];
  /*
    The params to use with the email template
   */
  params: TemplateType['params']; 
}

export const sendEmail = async <ParamsType extends Template>(params: EmailParams<ParamsType>): Promise<void> => {
  const { body, subject } = renderTemplate(params);

  const response = await resend.sendEmail({
    from: 'lukaskingsmail@gmail.com',
    to: params.to,
    html: body,
    subject,
  });

  const email = await resend.emails.get(response.id);
  email.
};
