type VerifyEmailTemplate = {
  template: 'verify-email';
  params: {
    userId: string,
    doiLink: string,
    expiresAt: string,
  }
}

type SigninTemplate = {
  template: 'verify-email';
  params: {
    ip: string;
    time: string;
  }
}

interface TemplateResult {
  subject: string;
  body: string;
}

export type Template = VerifyEmailTemplate | SigninTemplate;

export const renderTemplate = (template: Template): TemplateResult => {
  return {
    body: '',
    subject: '',
  }
};
