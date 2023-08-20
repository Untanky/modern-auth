import type { PostgresJsDatabase } from 'drizzle-orm/postgres-js';
import type { Preferences, PreferencesRepository } from './email.model';
import type * as schema from './schema';
import { preference } from './schema';
import { eq, sql } from 'drizzle-orm';

export class DrizzlePreferencesRepository implements PreferencesRepository {
    private readonly db: PostgresJsDatabase<typeof schema>;

    constructor(db: PostgresJsDatabase<typeof schema>) {
        this.db = db;
    }

    async findFirst(where?: Partial<Preferences> | undefined): Promise<Preferences> {
        const result = await this.db.query.preference.findFirst({
            // TODO: add better filters
            where: where && where.sub ? eq(preference.sub, where.sub) : sql`1 = 1`,
        });

        if (!result) {
            throw new Error('not found');
        }

        return {
            ...result,
            verifiedAt: result.verifiedAt || undefined,
        };
    }

    findMany(): Promise<Preferences[]> {
        throw new Error('Method not implemented.');
    }

    create(entity: Preferences): Promise<Preferences> {
        return this.db
            .insert(preference)
            .values(entity)
            .returning()
            .then(([result]) => ({
                ...result,
                verifiedAt: result.verifiedAt || undefined,
            }));
    }

    update(entity: Preferences): Promise<Preferences> {
        return this.db.update(preference)
            .set(entity)
            .where(eq(preference.sub, entity.sub))
            .returning()
            .then(([result]) => ({
                ...result,
                verifiedAt: result.verifiedAt || undefined,
            }));
    }

    delete(where: Partial<Preferences>): Promise<void> {
        return this.db.delete(preference)
            .where(where && where.sub ? eq(preference.sub, where.sub) : sql`1 = 1`)
            .then();
    }
}
