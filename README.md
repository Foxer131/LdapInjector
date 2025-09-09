# LdapInjector

LdapInjector is a Go-based tool designed for testing LDAP injection vulnerabilities. It allows you to perform blind LDAP injection attacks to extract data, such as passwords, by sending a series of requests to a target server.

-----

## How It Works

The tool works by iterating through a character set and sending requests to the target URL with a payload that attempts to guess the password one character at a time. It uses the server's response status code to determine whether the guessed character is correct.

The `LdapInjector` struct is the core of the tool, and it can be configured with different HTTP clients to perform the requests. Two clients are provided out of the box:

  * **`NetHttpBrute`**: Uses the standard Go `net/http` library.
  * **`FastHttpBrute`**: Uses the `fasthttp` library, which is designed for high-performance scenarios.

-----

## How to Use

To use the LdapInjector, you'll need to configure it with an HTTP client and then call the `BruteForce()` method. Here's a basic example of how to set it up:

```go
package main

import "fmt"

func main() {
    // Configure the HTTP client
    httpClient := NewHttpBrute(
        "POST",
        "http://your-target-url.com/login",
        "your-username",
        303, // The expected status code for a successful request
        map[string]string{
            "Content-Type": "application/x-www-form-urlencoded",
        },
    )

    // Create a new LdapInjector
    injector := NewLdapInjector(httpClient)

    // (Optional) Prune the character set to speed up the process
    injector.PruneCharset()

    // Run the brute-force attack
    password, err := injector.BruteForce()
    if err != nil {
        fmt.Println("Error:", err)
    }

    fmt.Println("Password found:", password)
}
```

-----

## Configuration

When creating a new HTTP client, you'll need to provide the following information:

  * **`verb`**: The HTTP method to use (e.g., `"POST"`, `"GET"`).
  * **`url`**: The URL of the target application.
  * **`username`**: The username to use in the LDAP injection payload.
  * **`expectedStatusCode`**: The HTTP status code that the server returns when a request is successful.
  * **`headers`**: A map of HTTP headers to include in the request.

-----

## Key Functions

Here are some of the key functions in the `LdapInjector`:

  * **`NewLdapInjector(client Injector)`**: Creates a new `LdapInjector` with the specified HTTP client.
  * **`PruneCharset()`**: Reduces the character set to only include characters that are valid in the target's password. This can significantly speed up the brute-force process.
  * **`TestCharacter(prefix string)`**: Tests a single character to see if it's the next valid character in the password.
  * **`BruteForce()`**: The main function that orchestrates the brute-force attack. It repeatedly calls `TestCharacter()` until the full password has been found.