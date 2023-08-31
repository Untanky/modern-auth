import type { PostgresJsDatabase } from 'drizzle-orm/postgres-js';
import type { ProfileRepository } from '../email.model';
import * as schema from '../schema'; 
import { eq, sql } from 'drizzle-orm';
import type { Profile } from '$lib/profile/model';

export class DrizzleProfileRepository implements ProfileRepository {
    readonly #db: PostgresJsDatabase<typeof schema>;

    constructor(db: PostgresJsDatabase<typeof schema>) {
        this.#db = db;
    }

    async findFirst(where?: Partial<Profile> | undefined): Promise<Profile> {
        const [result] = await this.#db
            .select()
            .from(schema.profile)
            .where(where && where.sub ? eq(schema.profile.sub, where.sub) : sql`1 = 1`)
            .limit(1);

        return result.data;
    }

    findMany(where?: Partial<Profile> | undefined): Promise<Profile[]> {
        throw new Error('Method not implemented.');
    }

    async create(entity: Profile): Promise<Profile> {
        const [result] = await this.#db
            .insert(schema.profile)
            .values({ sub: entity.sub, data: entity })
            .returning();

        return result.data;
    }

    async update(entity: Profile & { sub: string }): Promise<Profile> {
        const [result] = await this.#db
            .update(schema.profile)
            .set({ data: entity })
            .where(eq(schema.profile.sub, entity.sub))
            .returning();

        return result.data;
    }

    async delete(where: Partial<Profile>): Promise<void> {
        await this.#db
            .delete(schema.profile)
            .where(where && where.sub ? eq(schema.profile.sub, where.sub) : sql`1 = 1`);
    }
}
