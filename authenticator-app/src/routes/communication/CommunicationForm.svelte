<script lang="ts">
    import type { InsertPreferences, Preferences } from '$lib/preferences/model';
    import Input from '$lib/utils/Input.svelte';

    export let preferences: Preferences;

    export let onSave: (preferences: InsertPreferences) => void;

    let untouched = true;

    const onSubmit = (event: SubmitEvent): void => {
        const target = event.target as HTMLFormElement;
        onSave({
            emailAddress: (target[0] as HTMLInputElement).value,
            allowAccountReset: (target[1] as HTMLInputElement).checked,
            allowSessionNotification: (target[2] as HTMLInputElement).checked,
        });
    }
</script>

<form class="flex flex-col space-y-4" on:submit={onSubmit}>
    <section class="mt-2">
        <h2>
        Email Address
        {#if preferences.verified}
            <span class="pill pill-sm pill-green"> Verified </span>
        {:else}
            <span class="pill pill-red"> Verification required </span>
        {/if}
        </h2>
        <Input
            class="grow"
            id="email"
            label="Email"
            value={preferences.emailAddress}
            autocomplete="email"
            placeholder="abc@example.com"
            onInput={() => untouched = false}
        />
    </section>
    <section class="mt-2">
        <h2>Preferences</h2>
        <ul>
            <li>
                <input
                    id="allow-account-reset"
                    type="checkbox"
                    class="rounded-sm text-yellow-500"
                    checked={preferences.allowAccountReset}
                    on:change={() => untouched = false}
                />
                <label for="allow-account-reset">Use email to reset account</label>
            </li>
            <li>
                <input
                    id="allow-session-notification"
                    type="checkbox"
                    class="rounded-sm text-yellow-500"
                    checked={preferences.allowSessionNotification}
                    on:change={() => untouched = false}
                />
                <label for="allow-session-notification">Notify about new sessions</label>
            </li>
        </ul>
    </section>
    <button
        type="submit"
        class="btn btn-yellow self-end"
        disabled={untouched}
    >Update</button>
</form>