import { relations } from 'drizzle-orm';
import {
    boolean, pgEnum, pgSchema, timestamp, uuid, varchar,
} from 'drizzle-orm/pg-core';

export const emailSchema = pgSchema('email');

export const preference = emailSchema.table('preference', {
    sub: uuid('sub').primaryKey()
        .defaultRandom(),
    emailAddress: varchar('email_address').notNull(),
    verified: boolean('verified').notNull(),
    verifiedAt: timestamp('verified_at'),
    allowAccountReset: boolean('allow_account_reset').default(true)
        .notNull(),
    allowSessionNotification: boolean('allow_session_notification').default(false)
        .notNull(),
});

export const templates = pgEnum('templates_enum', [
    'verification',
    'accountReset',
    'sessionNotification',
]);

export const email = emailSchema.table('email', {
    id: uuid('id').primaryKey()
        .defaultRandom(),
    sub: uuid('sub').references(() => preference.sub, { onDelete: 'set null' }),
    sentAt: timestamp('sent_at').notNull(),
    template: templates('template').notNull(),
});

export const resendEmail = emailSchema.table('resend_email', {
    id: uuid('id').primaryKey()
        .references(() => email.id, { onDelete: 'cascade' }),
    resendId: varchar('resend_id').notNull()
        .unique(),
});

export const emailResendEmailRelation = relations(email, ({ one }) => ({
    resendEmail: one(resendEmail, {
        fields: [email.id],
        references: [resendEmail.id],
    }),
}));
