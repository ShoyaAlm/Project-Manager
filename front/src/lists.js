const lists = [
    {
        id: 1,
        name: 'doing',
        cards:[
            {   id: 1,
                name:'planning',
                members:['alex', 'josh', 'lucas', 'peter'],
                dates:['24th august', '21th september'],
                checklist:[
                    {
                        id: 1,
                        name: 'checklist 1',
                        item:[
                            {
                                id: 1,
                                name:'item 1',
                                dueDate:'24th september',
                                assignedTo:['josh', 'peter']
                            },
                            {
                                id: 2,
                                name:'item 2',
                                dueDate:'30th september',
                                assignedTo:['alex', 'lucas']
                            },

                        ]
                    },
                    {
                        id: 2,
                        name: 'checklist 2',
                        item:[
                            {
                                id: 1,
                                name:'item 1',
                                dueDate:'24th september',
                                assignedTo:['josh', 'peter']
                            },
                            {
                                id: 2,
                                name:'item 2',
                                dueDate:'30th september',
                                assignedTo:['alex', 'lucas']
                            },
                        ]
                    }
                ]
            },

            {   id: 2,
                name:'meeting',
                members:['alex', 'josh', 'lucas', 'peter'],
                dates:['24th august', '21th september'],
                checklist:[
                    {
                        id: 1,
                        name: 'checklist 1',
                        item:[
                            {
                                id: 1,
                                name:'item 1',
                                dueDate:'24th september',
                                assignedTo:['josh', 'peter']
                            },
                            {
                                id: 2,
                                name:'item 2',
                                dueDate:'30th september',
                                assignedTo:['alex', 'lucas']
                            },

                        ]
                    },
                    {
                        id: 2,
                        name: 'checklist 2',
                        item:[
                            {
                                id: 1,
                                name:'item 1',
                                dueDate:'24th september',
                                assignedTo:['josh', 'peter']
                            },
                            {
                                id: 2,
                                name:'item 2',
                                dueDate:'30th september',
                                assignedTo:['alex', 'lucas']
                            },
                        ]
                    }
                ]
            },
            
        ]
    },
    {
        id: 1,
        name: 'to do',
        cards:[
            {   id: 1,
                name:'planning',
                members:['alex', 'josh', 'lucas', 'peter'],
                dates:['24th august', '21th september'],
                checklist:[
                    {
                        id: 1,
                        name: 'checklist 1',
                        item:[
                            {
                                id: 1,
                                name:'item 1',
                                dueDate:'24th september',
                                assignedTo:['josh', 'peter']
                            },
                            {
                                id: 2,
                                name:'item 2',
                                dueDate:'30th september',
                                assignedTo:['alex', 'lucas']
                            },

                        ]
                    },
                    {
                        id: 2,
                        name: 'checklist 2',
                        item:[
                            {
                                id: 1,
                                name:'item 1',
                                dueDate:'24th september',
                                assignedTo:['josh', 'peter']
                            },
                            {   id: 2,
                                name:'item 2',
                                dueDate:'30th september',
                                assignedTo:['alex', 'lucas']
                            },
                        ]
                    }
                ]
            },

            {
                id: 2,
                name:'meeting',
                members:['alex', 'josh', 'lucas', 'peter'],
                dates:['24th august', '21th september'],
                checklist:[
                    {
                        id: 1,
                        name: 'checklist 1',
                        item:[
                            {
                                id: 1,
                                name:'item 1',
                                dueDate:'24th september',
                                assignedTo:['josh', 'peter']
                            },
                            {
                                id: 2,
                                name:'item 2',
                                dueDate:'30th september',
                                assignedTo:['alex', 'lucas']
                            },

                        ]
                    },
                    {
                        id: 2,
                        name: 'checklist 2',
                        item:[
                            {
                                id: 1,
                                name:'item 1',
                                dueDate:'24th september',
                                assignedTo:['josh', 'peter']
                            },
                            {
                                id: 2,
                                name:'item 2',
                                dueDate:'30th september',
                                assignedTo:['alex', 'lucas']
                            },
                        ]
                    }
                ]
            },
            
        ]
    }
]

export default lists