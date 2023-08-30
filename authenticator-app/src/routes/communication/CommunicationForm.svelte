<script lang="ts">
  import type { Preferences } from '$lib/preferences/model';
  import Input from '$lib/utils/Input.svelte';

  export let preferences: Preferences;

  export let onSave = (preferences: Preferences): void => {

  };
</script>

<section class="mt-2">
    <div>
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
            onInput={() => {}}
        />
    </div>
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
            />
            <label for="allow-account-reset">Use email to reset account</label>
        </li>
        <li>
            <input
                id="allow-session-notification"
                type="checkbox"
                class="rounded-sm text-yellow-500"
                checked={preferences.allowSessionNotification}
            />
            <label for="allow-session-notification">Notify about new sessions</label>
        </li>
    </ul>
</section>
<div class="flex justify-end mt-4 space-x-2">
    <button class="btn btn-yellow">Save</button>
</div>