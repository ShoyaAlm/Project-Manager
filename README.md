# Project-Manager

This is a Persian Project Management Web Application, and it is a clone from the Trello Web App.

---

**Getting Started:**

Once you have cloned the code, or downloaded a zip file, you can run the following commands in order : 

```bash
cd Project-Manager
```

*to run the back-end code, use this command:

```
go run main.go
```

*and to setup the front-end, simply run these commands:

```
cd front
npm install
npm start
```

Now you should have the app up and running. Easy as that :)

---

**Joining in:**


When you first run the program, you must either login or sign up if you haven't already. Click on the green button located at the top of the front-page of the web app.

<img src="https://github.com/ShoyaAlm/Project-Manager/assets/64843555/e6962f7a-4cc1-40df-9e3e-00f694e0d47e" alt="Login Image" width="240"> <img src="https://github.com/ShoyaAlm/Project-Manager/assets/64843555/4993df24-a2b2-4656-a2f3-0763e424c7c9" alt="Signup Image" width="300">




---

**Workspace:**

Once you sign in/login, you'll land on the 'Workspace' page. Here, all the boards you've created are listed above, while boards created by other users are listed below.

<img src="https://github.com/ShoyaAlm/Project-Manager/assets/64843555/ac6d1224-32b7-4eac-b7b5-71a1e7c413f2" alt="![Workspace]" width="640">

You can delete your boards by clicking on the red 'Remove' button.

---

**Boards, Lists, and Cards:**

Each 'Board' contains 'Lists', which in turn include 'Cards'. Cards are distinguished by different colors representing their labels. Additionally, each card has a starting date and a due date, editable by card members. Here's an example of a board :

<img src="https://github.com/ShoyaAlm/Project-Manager/assets/64843555/9a421d99-dce8-43a4-a0ce-2e098a775233" alt="![A Board eg]" width="640">

---

**Features:**

- **Drag & Drop:** Utilize the drag-and-drop feature seamlessly.

https://github.com/ShoyaAlm/Project-Manager/assets/64843555/51107e0f-2134-4d5e-a74f-1ccbfc2413de



- **Adding Lists:** Add new lists to your board by typing the name in the designated area.

<img src="https://github.com/ShoyaAlm/Project-Manager/assets/64843555/1d11ec4e-b635-414e-8aae-b1f73a139152" alt="![Add New List]" width="240">


- **Adding Cards:** Populate lists with cards effortlessly.

<img src="https://github.com/ShoyaAlm/Project-Manager/assets/64843555/eb67892e-8936-4758-a5a2-5eae297e001f" alt="![add-card-to-list]" width="240">



- **Table of Cards:** For every Board, you can see the cards inside the lists corresponding to that board shown like this :

<img src="https://github.com/ShoyaAlm/Project-Manager/assets/64843555/eddf681a-f8a8-4bd4-9be5-7a3a8549e26f" alt="![cards-table]" width="640">


*These infos are (from right to left) the name of the list they represent, the name of the card, their label, the members, and finally the due date*


---

**Card Details:**

When clicking on a card, you'll access its detailed view. Here, you can:

- Delete the card
- Manage card members
- Modify the description
- Add checklists
- Choose a label
- Change the dates


<img src="https://github.com/ShoyaAlm/Project-Manager/assets/64843555/81122704-523d-427e-9e7c-ffbf0f7e3e75" alt="![new-card]" width="500">


---

**Activity Section:**

The 'Activity' section logs all card-related activities, such as adding members, modifying checklists, adding labels, etc. You can also add comments for better collaboration.

Here's an example for the Activity section from another card :

<img src="https://github.com/ShoyaAlm/Project-Manager/assets/64843555/50a58c02-0e40-4319-b0e1-d290ddd75e16" alt="![activity-example]" width="300">

--- 


**Checklist Section:**

The section that is dedicated for making checklists looks like this :

<img src="https://github.com/ShoyaAlm/Project-Manager/assets/64843555/1e43e1e3-747c-4552-861d-66e6e7128635" alt="![checklist]" width="200">

*one checklists has items in it, while the other one(the one below) doesn't have any items.*






- **Items**

  We can also add items to our checklists  

  <img src="https://github.com/ShoyaAlm/Project-Manager/assets/64843555/aa6562d1-661c-4c84-9fb9-ebb74b25ef64" alt="![add-item]" width="300">


  Once we add a new item, we can change some parts of it. For instance, we can:
  
  - Tick the item (if it is done)
  - Assign the item to a card member
  - Modify the dates of an item
  - Delete item



  <img src="https://github.com/ShoyaAlm/Project-Manager/assets/64843555/a2b06cba-f467-4f6b-881e-2156f4828466" alt="![item]" width="300">

---



This web app was part of my 'Bachelor Project' and I took a lot of influence from similar web apps such as Trello, Basecamp and Notion.
