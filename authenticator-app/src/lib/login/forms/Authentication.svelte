<script lang="ts">
  import Input from '../../utils/Input.svelte';
import RadioBox from '../../utils/RadioBox.svelte';
import RadioGroup from '../../utils/RadioGroup.svelte';

  export let submit: () => void;
  export let userId: string;

  let method: 'public-key' | 'password' = 'public-key';

  let password = '';

  const onSubmit = (event: SubmitEvent) => {
      event.preventDefault();

      submit();
  };

  const onChangePassword = (event: Event): void => {
      const input = event.target as HTMLInputElement;

      password = input.value;
      return;
  };
</script>

<form class="flex flex-col flex-1 space-y-2" on:submit={onSubmit}>
  <h2 class="text-xl">
    Authentication
  </h2>
  <Input
    label="User ID:"
    id="user-id"
    autocomplete="username"
    value={userId}
    onInput={() => {}}
    disabled
  />
  <p>
    Please select how to authenticate:
  </p>
  <RadioGroup class="space-y-4">
    <RadioBox
      open={method === 'public-key'}
      click={() => {
          method = 'public-key';
      }}
    >
      <h3 class="text-lg">Biometric or physical authentication</h3>
      <div class="{method === 'public-key' ? '' : 'hidden'} flex flex-col">
        <p>
          When you click on authenticate a system dialog will open and ask you
          to authenticate with your biometric data or a physical hardware token.
          Please prepare for the method chosen when setting up this device.
        </p>
        <button type="submit" class="self-end mt-2 btn btn-yellow">
          Authenticate
        </button>
      </div>
    </RadioBox>
    <RadioBox
      open={method === 'password'}
      click={() => {
          method = 'password';
      }}
    >
      <h3 class="text-lg">Password authentication</h3>
      <div class="{method === 'password' ? '' : 'hidden'} flex flex-col">
        <p>
          Please enter your password:
        </p>
        <Input
          label="Password:"
          id="password"
          type="password"
          autocomplete="current-password"
          value={password}
          onInput={onChangePassword}
        />
        <button type="submit" class="self-end mt-2 btn btn-yellow">
          Continue
        </button>
      </div>
    </RadioBox>
  </RadioGroup>
</form>
