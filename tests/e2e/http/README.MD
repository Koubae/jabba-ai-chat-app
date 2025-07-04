Testing e2e -- http
===================



#### Admin Access-Token

```bash
access_token_admin=$(curl -s --request POST  localhost:20000/api/v1/auth/login -H "Content-Type: application/json" -d '{"username": "admin","password": "admin","application_id": "ai-bot"}' | jq -r '.access_token'); \
  echo $access_token_admin;
```

#### Sign UP

```bash
user_id=$(curl -s --request POST localhost:20000/api/v1/auth/signup -H "Content-Type: application/json" -d '{"username": "user_1", "password": "pass","application_id": "ai-bot"}' | jq -r '.id'); \
  echo $user_id;
  
curl -s --request POST localhost:20000/api/v1/auth/signup -H "Content-Type: application/json" -d '{"username": "user_1", "password": "pass","application_id": "ai-bot"}';echo;
```

#### Login

```bash
access_token=$(curl -s --request POST  localhost:20000/api/v1/auth/login -H "Content-Type: application/json" -d '{"username": "user_1","password": "pass","application_id": "ai-bot"}' | jq -r '.access_token'); \
  echo $access_token; \
  curl localhost:20000/api/v1/account/1?full_profile=false -H "Authorization: Bearer $access_token"
```

#### SignUp + Login

```bash
user_id=$(curl -s --request POST localhost:20000/api/v1/auth/signup -H "Content-Type: application/json" -d '{"username": "user_1", "password": "pass","application_id": "ai-bot"}' | jq -r '.user_id'); \
  echo $user_id; echo -e "\n"; \
  access_token=$(curl -s --request POST  localhost:20000/api/v1/auth/login -H "Content-Type: application/json" -d '{"username": "user_1","password": "pass","application_id": "ai-bot"}' | jq -r '.access_token'); \
  echo $access_token; echo -e "\n"; \
  curl localhost:20000/api/v1/account/$user_id?full_profile=false -H "Authorization: Bearer $access_token"
```


```bash
# signup
curl -s --request POST localhost:20000/api/v1/auth/signup -H "Content-Type: application/json" -d '{"username": "user_1", "password": "pass","application_id": "ai-bot"}';echo;
# login
curl -s --request POST  localhost:20000/api/v1/auth/login -H "Content-Type: application/json" -d '{"username": "user_1","password": "pass","application_id": "ai-bot"}'
```