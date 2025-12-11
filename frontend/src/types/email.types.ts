export type EmailProvider = "mailto" | "resend" | "smtp";

export interface InvoiceEmailSettings {
  provider: EmailProvider;
  from?: string;
  replyTo?: string;
  subjectTemplate?: string;
  bodyTemplate?: string;
  resendApiKey?: string;
  smtpHost?: string;
  smtpPort?: number;
  smtpUsername?: string;
  smtpPassword?: string;
  smtpUseTLS?: boolean;
  signature?: string;
}
