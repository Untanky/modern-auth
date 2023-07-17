<script lang="ts">
  import { initiateAuthentication } from "./Authentication";
  import CreateCredential from "./CreateCredential.svelte";
  import GetCredential from "./GetCredential.svelte";
  import UserIdForm from "./UserIdForm.svelte";

  type State = 'userId' | 'createCredential' | 'getCredential' | 'success';
  
  let state: State = 'userId';
  let loading = false;
  let error: Error | null = null;
  let userId: string | null = null;

  const onInitiateAuthentication = (): void => {
    const inputElement = document.getElementById('user-id') as HTMLInputElement;
    const inputValue = inputElement.value;

    userId = inputValue;
    loading = true;

    initiateAuthentication(inputValue)
      .then(() => {
        state = 'createCredential';
      })
      .catch((err) => {
        error = err;
      })
      .finally(() => {
        loading = false;
      });
  };

  const onCreateCredential = (): void => {
    console.log('create credential');
  };

  const onGetCredential = (): void => {
    console.log('create credential');
  };

</script>

<div class="flex w-screen h-screen justify-center items-center">
  <main class="w-[420px] card">
    <h1 class="text-stone-950 dark:text-stone-50 text-xl font-medium">
      Login
    </h1>
    {#if state === 'userId'}
    <UserIdForm submit={onInitiateAuthentication} />
    {:else if state === 'createCredential'}
    <CreateCredential userId={userId} submit={onCreateCredential} />
    {:else if state === 'getCredential'}
    <GetCredential userId={userId} submit={onGetCredential} />
    {/if}
  </main>
</div>
