// export const userData = [
//   {
//     id: 1,
//     messages: [
//       {
//         id: 1,
//         name: "Jane Doe",
//         message: "Hey, Jakob",
//       },
//       {
//         id: 2,
//         name: "Jakob Hoeg",
//         message: "Hey!",
//       },
//       {
//         id: 3,
//         name: "Jane Doe",
//         message: "How are you?",
//       },
//       {
//         id: 4,
//         name: "Jakob Hoeg",
//         message: "I am good, you?",
//       },
//       {
//         id: 5,
//         name: "Jane Doe",
//         message: "I am good too!",
//       },
//       {
//         id: 6,
//         name: "Jakob Hoeg",
//         message: "That is good to hear!",
//       },
//       {
//         id: 7,
//         name: "Jane Doe",
//         message: "How has your day been so far?",
//       },
//       {
//         id: 8,
//         name: "Jakob Hoeg",
//         message:
//           "It has been good. I went for a run this morning and then had a nice breakfast. How about you?",
//       },
//       {
//         id: 9,
//         name: "Jane Doe",
//         message: "I had a relaxing day. Just catching up on some reading.",
//       },
//     ],
//     name: "Jane Doe",
//   },
//   {
//     id: 2,
//     messages: [
//       {
//         id: 1,
//         name: "John Doe",
//         message: "Hey, Jakob",
//       },
//       {
//         id: 2,
//         name: "Jakob Hoeg",
//         message: "Hey!",
//       },
//       {
//         id: 3,
//         name: "John Doe",
//         message: "How are you?",
//       },
//       {
//         id: 4,
//         name: "Jakob Hoeg",
//         message: "I am good, you?",
//       },
//       {
//         id: 5,
//         name: "John Doe",
//         message: "I am good too!",
//       },
//       {
//         id: 6,
//         name: "Jakob Hoeg",
//         message: "That is good to hear!",
//       },
//       {
//         id: 7,
//         name: "John Doe",
//         message: "How has your day been so far?",
//       },
//       {
//         id: 8,
//         name: "Jakob Hoeg",
//         message:
//           "It has been good. I went for a run this morning and then had a nice breakfast. How about you?",
//       },
//       {
//         id: 9,
//         name: "John Doe",
//         message: "I had a relaxing day. Just catching up on some reading.",
//       },
//     ],
//     name: "John Doe",
//   },
//   {
//     id: 3,
//     messages: [
//       {
//         id: 1,
//         name: "Elizabeth Smith",
//         message: "Hey, Jakob",
//       },
//       {
//         id: 2,
//         name: "Jakob Hoeg",
//         message: "Hey!",
//       },
//       {
//         id: 3,
//         name: "Elizabeth Smith",
//         message: "How are you?",
//       },
//       {
//         id: 4,
//         name: "Jakob Hoeg",
//         message: "I am good, you?",
//       },
//     ],
//     name: "Elizabeth Smith",
//   },
//   {
//     id: 4,
//     messages: [
//       {
//         id: 1,
//         name: "John Smith",
//         message: "Hey, Jakob",
//       },
//       {
//         id: 2,
//         name: "Jakob Hoeg",
//         message: "Hey!",
//       },
//       {
//         id: 3,
//         name: "John Smith",
//         message: "How are you?",
//       },
//       {
//         id: 4,
//         name: "Jakob Hoeg",
//         message: "I am good, you?",
//       },
//       {
//         id: 5,
//         name: "John Smith",
//         message: "I am good too!",
//       },
//       {
//         id: 6,
//         name: "Jakob Hoeg",
//         message: "That is good to hear!",
//       },
//       {
//         id: 7,
//         name: "John Smith",
//         message: "How has your day been so far?",
//       },
//     ],
//     name: "John Smith",
//   },
// ];

export type Friend = {
  id: number;
  username: string;
  messages: Message[];
  room_id: number;
};

export const loggedInUserData = {
  id: 5,
  name: "Jakob Hoeg",
};

export type LoggedInUserData = typeof loggedInUserData;

export interface Message {
  id: number;
  user_id: number;
  message: string;
}
