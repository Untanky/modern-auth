import type { Template } from './email.model';

type TemplateResult = {
  body: string;
  subject: string;
};

export const renderTemplate = (template: Template): TemplateResult => {  
  return {
    body: '',
    subject: '',
  }
};
