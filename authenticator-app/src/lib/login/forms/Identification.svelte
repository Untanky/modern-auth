<script lang="ts">
  import Input from '../../utils/Input.svelte';

  export let submit: (userId: string) => void;

  let userId = '';

  $: canSubmit = userId.length > 0;

  const onClick = (event: SubmitEvent) => {
      event.preventDefault();

      const formElement = event.target as HTMLFormElement;
      const inputElement = formElement[0] as HTMLInputElement;

      submit(inputElement.value);
  };

  const onChangeUserId = (event: Event): void => {
      const input = event.target as HTMLInputElement;

      userId = input.value;
  };
</script>

<form class="flex flex-col flex-1" on:submit={onClick}>
  <h2 class="text-xl">
    Identification
  </h2>
  <p class="mt-2">
    To authenticate, we need to know who you are. Please enter your user id.
  </p>
  <Input
    label="User ID:"
    value={userId}
    onInput={onChangeUserId}
    id="user-id"
    type="text"
    placeholder="Your user id"
    autocomplete="username"
  />
  <button type="submit" class="self-end mt-4 btn btn-yellow" disabled={!canSubmit}>
    Continue
  </button>
</form>
