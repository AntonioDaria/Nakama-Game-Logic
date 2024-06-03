# Game Logic Implementation with Nakama

This project implements a game logic solution using Nakama, an open-source game backend. 
The solution includes an RPC function to read a file, save its content to the database, and verify the content's hash. 
Additionally, a health check RPC function is provided.


## What's been accomplished ✅

Implemented an RPC function fileHandler that:
- Accepts payload parameters (type, version, hash).
- Reads a file from disk based on the type and version.
- Calculates the file's hash and compares it to the provided hash.
- Saves the file content to Nakama's storage engine if the hashes match.
- Returns the file content and metadata in the response.
- Implemented a health check RPC function healthcheck.
- Added unit tests for the RpcFileHandler and storage save functionality.


## Getting started

To run it locally 

### Prerequisites

- Docker
- Docker Compose
- Make

### Building and Running the Application

1. **Clone the repository:**
```sh
git clone <repository-url>
cd <repository-directory>
```
   
2. **Start the application:**
```sh
make start
```


## Interacting with the API
```
You can interact with the RPC functions using Nakama's Developer Console.

Using Nakama's Developer Console
Access the Developer Console:

Navigate to http://localhost:7351.
Invoke the filehandler RPC Function:


On API Explore Select the fileHandler endpoint from the dropdown menu and set the user ID to `00000000-0000-0000-0000-000000000000`

Use the following example payload to test the function:

{
    "type": "core",
    "version": "1.0.0",
    "hash": "85a11a7d406be88cb7e0b2c68b8134fc" //Valid hash
}

You should observe the following behavior:
- If matching hash is provided the content will have data
- If no matching hash is provided the content will be null
- If an empty payload is provided the defaults will be used
- If no payload at all is provided the request will error with bad input


Invoke the healthCheck RPC Function:

Select the healthCheck endpoint from the dropdown menu.
No payload is needed for the health check.
```

### Verifying Data Saved in Nakama
To verify that the data has been saved to Nakama's storage, use the ReadStorageObjects endpoint:

1-Access the Developer Console:

    Navigate to http://localhost:7351.

2- Select the ReadStorageObjects Endpoint:

    Choose ReadStorageObjects from the dropdown menu.

3- Provide the Payload:

    Use the following JSON payload to read the stored data:

```
{
    "object_ids": [
        {
            "collection": "core",
            "key": "1.0.0",
            "user_id": ""
        }
    ]
}
```

4- Send the Request:
```
Click "Send Request" and observe the response to verify the stored data.
```

---

### Technologies used
- Golang: The primary language for the game logic implementation.
- Nakama: Open-source game backend for the server-side logic.
- Docker: Containerization for easy setup and deployment.
- Docker Compose: Tool for defining and running multi-container Docker applications.
- Testify: Assertion library for unit testing in Go.

---

###  Testing
All current tests are passing. Tests use Gomock and stretchr/testify.
This will run all the tests in the project.


3. **To Run the tests:**
```sh
make test
```

4. **Stop the application:**
```sh
make stop
```

---


### Explanation of the Solution

The RpcFileHandler function is designed to handle incoming RPC (Remote Procedure Call) requests in a Nakama server context. This function processes a JSON payload, reads and compacts a corresponding JSON file, calculates its hash, and optionally saves the compacted content to Nakama's storage engine if certain conditions are met. Here is a breakdown of the key steps involved:

1. **Initialization and Default Settings:**
The function initializes a Payload structure with default values for Type, Version, and Hash. It logs the received payload for debugging purposes.

2. **Payload Validation:**
If the payload is empty, the function logs an error and returns an error response. It then attempts to unmarshal the payload into the Payload structure, logging an error and returning if unmarshalling fails.

3. **Use of Default Values:**
The useDefaults function is called to ensure that the Payload structure has default values if any fields are not provided in the payload.

4. **File Reading and Compacting:**
The readAndCompactFile function constructs the file path based on the payload's Type and Version. It reads the file's content and compacts the JSON content to remove unnecessary whitespace. If the file cannot be read or the JSON compaction fails, the function logs an error and returns.

5. **Hash Calculation:**
The calculateHash function calculates the MD5 hash of the compacted file content and logs the calculated hash. While MD5 has known vulnerabilities and is not recommended for cryptographic security,
its use in our scenario is justified given the trade-offs between performance and security requirements for simple integrity verification, as it is a relatively simple and fast hashing algorithm, which makes it efficient for quick hash generation.
I think the choice of MD5 it is adequate for basic integrity checks where the risk of deliberate hash collisions is minimal.

6. **Response Preparation and Storage:**
A Response structure is prepared with the Type, Version, and calculated hash. If the provided hash matches the calculated hash, the function saves the compacted content to Nakama's storage engine using the SaveToStorageEngine function.
The response is then marshaled into JSON and returned.


---


### Thoughts on the Task

1. **Robustness and Error Handling:**
The function includes extensive error handling, logging errors at various stages to aid in debugging. This is crucial in a production environment to diagnose issues quickly and effectively.

2. **Modularity:**
The task is broken down into smaller functions (useDefaults, readAndCompactFile, calculateHash, SaveToStorageEngine), promoting modularity and reusability. This makes the code easier to maintain, test and extend.

3. **Defaults and Flexibility:**
By providing default values and allowing payload overrides, the function is flexible and can handle a variety of input scenarios. This makes it more resilient to changes in input data.

4. **Logging:**
The use of logging at different stages of the process helps in tracking the flow of data and identifying points of failure.


---

### Ideas for Improvement

1. **Enhanced Validation:**
Additional validation on the payload fields (e.g., ensuring Type and Version conform to expected formats) could improve the robustness of the function.

2. **Configurable Paths and Settings:**
Instead of hardcoding file paths and default values, consider making these configurable through environment variables or configuration files. This would make the function more adaptable to different environments.

3. **Api Documentation:**
If I had more time, I could have also integrated Swagger documentation for the API. This would provide several benefits:

    - API Documentation: Swagger would automatically generate interactive and detailed documentation of the API endpoints, making it easier for developers to understand how to interact with the RPC handler.

    - Improved Developer Experience: With Swagger's interactive interface, developers can test the API endpoints directly from the documentation, reducing the learning curve and increasing productivity.

4. **Modularity**
On the topic of Modularity I could have probably refactored the rpc handler a bit more by extracting the majority of the logic out to a service function and leave the handler just to deal wit the http response codes.
I have not put a too much thoughts into this last point though, but it is something that I would normally have done if I was designing an API endpoint using an http framework like Fiber for example.
I think for the purpose of this exercise breaking done the handler into smaller function to enhance readability and testability was a simpler choice, but definitely something to explore further.

5. **Structure**
With more time It would have been better to move the smaller functions that make up the rpc handler, into separate files as well as moving their tests into a separate test files.
Again this would have increased Readability and improve maintainability.



