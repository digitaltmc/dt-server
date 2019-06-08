import { ApolloClient } from 'apollo-client';
import { HttpLink } from 'apollo-link-http';
import { InMemoryCache } from 'apollo-cache-inmemory';
import fetch from 'node-fetch';

import gql from 'graphql-tag';

const client = new ApolloClient({
    ssrMode: false,
    link: new HttpLink({
      uri: 'http://localhost:58080/graphql',
      fetch,
    }),
    cache: new InMemoryCache(),
});

let token = null;

const queryFunc = (input) => {
  return client.query({
    query: input,
  })
    .then(data => {
      if (data.data && data.data.meeting) {
        console.log("Meeting: ", data.data.meeting)
      } else {
        console.log(data);
      }
    })
    .catch(error => console.error(error))
    .finally(() => console.log("----------"));
};

const mutationFunc = (input) => {
  return client.mutate({
    mutation: input,
  })
    .then(data => {
      if (data.data && data.data.meeting) {
        console.log("Meeting: ", data.data.meeting)
      } else {
        console.log(data);
      }
      if (data.data && data.data.login) {
        token = data.data.login
        console.log("Set token:", token);
      }
    })
    .catch(error => console.error(error))
    .finally(() => console.log("----------"));
};

const createPromise = (type, input) => {
  if (typeof(input) === "function") input=input();
  if (type === "query") {
    return queryFunc(input);
  }
  return mutationFunc(input);
}

//---------- gql

const hello = gql`
  {
    hello
  }
`;

const deleteUser = gql`
mutation {
  deleteUser(username: "owen")
}
`

const registerUser = gql`
mutation {
  register(person:{name:"owen",password:"pwd",email:"aaa@sap.com",mobile:"888888"})
}
`

const loginUser = gql`
{
  login(user:"owen",password:"pwd")
}
`

const unbook1 = () => gql`
mutation {
  unbook(token: "${token}", date: "2011-01-01", roleName: "ie")
}
`
const unbook2 = () => gql`
mutation {
  unbook(token: "${token}", date: "2011-01-01", roleName: "ie2")
}
`

const book1 = () => gql`
mutation {
  book(token: "${token}", date: "2011-01-01", roleName: "ie", title: "ieTitle1")
}
`
const book2 = () => gql`
mutation {
  book(token: "${token}", date: "2011-01-01", roleName: "ie2", title: "ieTitle2")
}
`

const meeting = gql`
  {
    meeting(date: "2011-01-01") {
      date,
      agenda {
        roleName
        member {
          name
        }
       title 
      }
    }
  }
`

//----------

// Scenario 0: test

queryFunc(hello)

// Scenario 1:
// Delete user, register, login, book 1, book 2, update 1, list meeting.

// console.log(
// [
// mutationFunc(deleteUser),
// mutationFunc(registerUser),
// mutationFunc(loginUser),
// queryFunc(meeting)
// ]
// );

const scenario1 = [
['mutation', deleteUser],
['mutation', registerUser],
['mutation', loginUser],
['mutation', unbook1],
['mutation', book1],
['mutation', book2],
['mutation', unbook2],
['query', meeting]
];

scenario1.reduce( async (previousPromise, nextInput) => {
  await previousPromise;
  return createPromise(nextInput[0], nextInput[1]);
}, Promise.resolve());

