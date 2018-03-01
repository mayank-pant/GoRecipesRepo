### RECIPES 

PREREQUISITES :-

- docker
- docker-compose
- internet to initially download the docker images

HOW TO RUN :-

1. Go to the root folder where the code is present/cloned.
2. type -> ' docker-compose up'
3. Using POSTMAN can do queries by using localhost/<pathname>.
   For POST requests, under Raw, send the POST data as necessary.
4. Protected paths need a header as in "Authorization":"Bearer <token>"

LANGUAGE USED :-
Golang

HIGHLIGHTS :-

1. One liner setup. Code internally creates the necessary DB schema.
2. Caching is used to make the actions faster.
3. ORM has been used to virtualize the database actions.
4. Pagination is supported and you can customize pagination action.
5. Testing is done using Go Test.
6. gofmt formatting of code and includes comments.
7. JWT has been used to do authorization checks by providing token.
8. Protected API paths are present which do authorization checks as well to ensurw only the creator of a recipe can modify it.

API DETAILS :-

/signup - POST

- {"username":<stringvalue>,"password":<stringvalue>}

Response : {"success"}

/authenticate - POST 

- {"username":<stringvalue>,"password":<stringvalue>}

Response : {"token":<stringvalue>}

/search/recipes - GET

- Optional params are 
  - limit which tells the maximum search results to be shown
  - page  which tells the paginated page that is necessary

  LIMIT and PAGE both need to be given if results are to be PAGINATED.

- Use one of the following search params in the query string -
    - name
    - vegetarian
    - preptime
    - difficulty

Ex: - localhost/search/recipes?vegetarian=false&&limit=2&&page=2

Pagination is supported for Listing recipes as well similar to what is seen above.

Ex:- localhost/recipes?limit=1&page=2


THANKS.

SANDEEP BHAT
Any queries and you can contact sandyethadka@gmail.com