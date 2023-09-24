variable VERSION {
    default = "0.1-beta"
}

variable IS_LATEST {
    default = true
}

variable REF {
    default = ""
}

group default {
    targets = [
        "authenticator-app"
    ]
}

target authenticator-app {
    dockerfile = "Dockerfile"
    context = "authenticator-app"
    labels = {
        "org.opencontainers.image.name" = "modern-auth/authenticator-app"
        "org.opencontainers.image.description" = "App containing the frontend, email and profile service of ModernAuth"
        "org.opencontainers.image.url" = "https://github.com/Untanky/modern-auth"
        "org.opencontainers.image.source" = "https://github.com/Untanky/modern-auth"
        "org.opencontainers.image.revision" = "${REF}"
        "org.opencontainers.image.version=" = "${VERSION}"
    }
    tags = [
        // "untanky/authenticator-app:${VERSION}",
        IS_LATEST ? "ghcr.io/untanky/authenticator-app:latest" : "",
        "ghcr.io/untanky/authenticator-app:${VERSION}",
    ]
    cache-from = [
        "type=gha"
    ]
    cache-to = [
        "type=gha"
    ]
    platforms = [
        "linux/amd64"
    ]
}
