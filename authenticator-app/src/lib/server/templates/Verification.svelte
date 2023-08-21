<script lang="ts">
    import { Button, Container, Html, Link, Text } from 'svelte-email';
    export let userName: string = '';
    export let firstName: string = '';
    export let lastName: string = '';
    export let doiLink: string;
</script>

<Html lang="en">
    <Container>
        <Text>
            Hello {userName},
        </Text>
        <Text>
            your email has been used to register for ModernAuth. To verify you email please click on the button.
        </Text>
        <Button style={{ backgroundColor: 'rgb(250, 204, 21)', padding: '8px 16px', borderRadius: '16px' }} href={doiLink}>
            Verify
        </Button>
        <Text>
            If you can click the button, copy the following URL into your browser:
            <Link href={doiLink}>
                {doiLink}
            </Link>
        </Text>
    </Container>
</Html>