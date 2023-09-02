<script lang="ts">
    import type { Profile } from '$lib/profile/model';
  import Input from '$lib/utils/Input.svelte';
    
    export let profile: Profile;

    export let onSave: (profile: Profile) => void;

    let untouched = true;

    const onSubmit = (event: SubmitEvent): void => {
        const target = event.target as HTMLFormElement;
        console.log(target['firstname']);
        onSave({
            sub: '',
            firstName: target['firstname'].value,
            lastName: target['lastname'].value,
            userName: target['username'].value,
            emailAddress: target['email'].value,
            street: target['street'].value,
            postalCode: target['zip'].value,
            city: target['city'].value,
            country: target['country'].value,
        });
    }
</script>

<form class="flex flex-col pt-4 space-y-4" on:submit={onSubmit}>
    <div class="grid grid-cols-2 gap-2">
        <h3 class="text-lg col-span-2">Personal information</h3>
        <Input
            type="text"
            value={profile.firstName || ''}
            placeholder="Nomen"
            autocomplete="given-name"
            onInput={() => untouched = false}
            id="firstname"
            label="firstname"
        />
        <Input
            type="text"
            value={profile.lastName || ''}
            placeholder="Nescio"
            autocomplete="family-name"
            onInput={() => untouched = false}
            id="lastname"
            label="lastname"
        />
        <Input
            class="col-span-2"
            type="text"
            value={profile.userName || ''}
            placeholder="Your Username"
            autocomplete="username"
            onInput={() => untouched = false}
            id="username"
            label="username"
        />
        <Input
            class="col-span-2"
            type="text"
            value={profile.emailAddress || ''}
            placeholder="nomen.nescio@email.com"
            autocomplete="email"
            onInput={() => untouched = false}
            id="email"
            label="email"
        />
        <h3 class="text-lg">Address</h3>
        <Input
            class="col-span-2"
            type="text"
            value={profile.street || ''}
            placeholder="1234 Example Street"
            autocomplete="street-address"
            onInput={() => untouched = false}
            id="street"
            label="street"
        />
        <Input
            type="text"
            value={profile.postalCode || ''}
            placeholder="12345"
            autocomplete="postal-code"
            onInput={() => untouched = false}
            id="zip"
            label="zip"
        />
        <Input
            type="text"
            value={profile.city || ''}
            placeholder="Example City"
            autocomplete="address-level2"
            onInput={() => untouched = false}
            id="city"
            label="city"
        />
        <Input
            class="col-span-2"
            type="text"
            value={profile.country || ''}
            placeholder="Example Country"
            autocomplete="country-name"
            onInput={() => untouched = false}
            id="country"
            label="country"
        />
    </div>
    <button
        type="submit"
        class="btn btn-yellow self-end"
        disabled={untouched}
    >
        Update
    </button>
</form>