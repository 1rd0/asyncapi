   curl -X POST http://localhost:4000/auth/refreshTokens \
   -H "Content-Type: application/json" \
   -d '{"refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoiUmVmcmVzaCIsImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3Q6NDAwMCIsInN1YiI6Ijg3YzZmZmI3LTgzMzktNDEzOS04OGUwLWUxOTJhYTNiZTE2YiIsImV4cCI6MTc0MTA5MzA3NywiaWF0IjoxNzM4NTAxMDc3fQ.AqbZ093uuiAYgdK0Ovm3tfLsRz7sg3RG6C2xDPT0YuA"}'