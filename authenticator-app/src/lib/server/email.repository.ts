import type { PostgresJsDatabase } from 'drizzle-orm/postgres-js';
import type {  } from 'drizzle-orm/pg-core';
import type { Email, EmailRepository, Template } from "./email.model";
import { email, resendEmail } from './schema';
import { eq, sql } from 'drizzle-orm';



export class DrizzleEmailRepository implements EmailRepository {
  private readonly db;
  
  constructor(db: PostgresJsDatabase) {
    this.db = db;
  }

  findFirst(where?: Partial<Email> | undefined): Promise<Email> {
    return this.db
      .select()
      .from(email)
      .innerJoin(resendEmail, eq(email.id, resendEmail.id))
      // TODO: add better filters
      .where(where && where.id ? eq(email.id, where.id) : sql`1 = 1`)
      .limit(1)
      .then(([result]) => ({
        id: result.email.id,
        sub: result.email.sub || 'deleted',
        sentAt: result.email.sentAt,
        template: {
          type: result.email.template,
          sub: result.email.sub || 'deleted',
        } as Template,
        deliveryMethod: 'resend',
        resendId: result.resend_email?.resendId || ''
      }));
  }

  findMany(where?: Partial<Email> | undefined): Promise<Email[]> {
    return this.db
      .select()
      .from(email)
      .innerJoin(resendEmail, eq(email.id, resendEmail.id))
      // TODO: add better filters
      .where(where && where.id ? eq(email.id, where.id) : sql`1 = 1`)
      .then((results) => results.map((result) => ({
        id: result.email.id,
        sub: result.email.sub || 'deleted',
        sentAt: result.email.sentAt,
        template: {
          type: result.email.template,
          sub: result.email.sub || 'deleted',
        } as Template,
        deliveryMethod: 'resend',
        resendId: result.resend_email?.resendId || ''
      })));
  }

  create(model: Required<Email>): Promise<Email> {
    return this.db.transaction(async (tx): Promise<Email> => {
      const [createdEmail] = await tx.insert(email).values({
        sentAt: model.sentAt,
        template: model.template.type,
        sub: model.sub,
      }).returning({ insertedId: email.id });
      await tx.insert(resendEmail).values({
        id: createdEmail.insertedId,
        resendId: model.resendId,
      });
      return {
        ...model,
        id: createdEmail.insertedId,
      };
    })
  }

  update(entity: Email): Promise<Email> {
    throw new Error("Method not implemented.");
  }

  async delete(where: Partial<Email>): Promise<void> {
    await this.db
      .delete(email)
      .where(where && where.id ? eq(email.id, where.id) : sql`1 = 1`);
  }
}