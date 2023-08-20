import { eq, sql } from 'drizzle-orm';
import type { PostgresJsDatabase } from 'drizzle-orm/postgres-js';
import type {
    Email, EmailRepository, Template,
} from './email.model';
import type * as schema from './schema';
import { email, resendEmail } from './schema';

export class DrizzleEmailRepository implements EmailRepository {
    private readonly db: PostgresJsDatabase<typeof schema>;

    constructor(db: PostgresJsDatabase<typeof schema>) {
        this.db = db;
    }

    findFirst(where?: Partial<Email> | undefined): Promise<Email> {
        return this.db.query.email.findFirst({
            // TODO: add better filters
            where: where && where.id ? eq(email.id, where.id) : sql`1 = 1`,
            with: { resendEmail: true },
        }).then((result) => {
            if (!result) {
                throw new Error('not found');
            }

            return {
                id: result.id,
                sub: result.sub || '',
                template: {
                    type: result.template,
                    sub: result.sub || '',
                } as Template,
                sentAt: result.sentAt,
                deliveryMethod: 'resend',
                resendId: result.resendEmail.resendId,
            };
        });
    }

    findMany(where?: Partial<Email> | undefined): Promise<Email[]> {
        return this.db.query.email.findMany({
            // TODO: add better filters
            where: where && where.id ? eq(email.id, where.id) : sql`1 = 1`,
            with: { resendEmail: true },
        }).then((results) => results.map((result) => ({
            id: result.id,
            sub: result.sub || '',
            template: {
                type: result.template,
                sub: result.sub || '',
            } as Template,
            sentAt: result.sentAt,
            deliveryMethod: 'resend',
            resendId: result.resendEmail.resendId,
        })),
        );
    }

    create(model: Required<Email>): Promise<Email> {
        return this.db.transaction(async (tx): Promise<Email> => {
            const [createdEmail] = await tx.insert(email).values({
                sentAt: model.sentAt,
                template: model.template.type,
                sub: model.sub,
            })
                .returning({ insertedId: email.id });
            await tx.insert(resendEmail).values({
                id: createdEmail.insertedId,
                resendId: model.resendId,
            });
            return {
                ...model,
                id: createdEmail.insertedId,
            };
        });
    }

    update(): Promise<Email> {
        throw new Error('Method not implemented.');
    }

    async delete(where: Partial<Email>): Promise<void> {
        await this.db
            .delete(email)
            .where(where && where.id ? eq(email.id, where.id) : sql`1 = 1`);
    }
}
