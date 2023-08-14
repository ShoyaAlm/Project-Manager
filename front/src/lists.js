const lists = [
    {
        id: 1,
        name: 'doing',
        cards:[
            {   id: 1,
                name:'planning',
                members:['alex', 'josh', 'lucas', 'peter'],
                dates:['24th august', '21th september'],
                description:'just some description',
                checklists:[
                    {
                        id: 1,
                        name: 'checklist 1',
                        items:[
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
                        items:[
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
                description:'another description',
                checklists:[
                    {
                        id: 1,
                        name: 'checklist 1',
                        items:[
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
                        items:[
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
        id: 2,
        name: 'to do',
        cards:[
            {   id: 1,
                name:'working',
                members:['alex', 'josh', 'lucas', 'peter'],
                dates:['24th august', '21th september'],
                description:'just some description',
                checklists:[
                    {
                        id: 1,
                        name: 'checklist 1',
                        items:[
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
                        items:[
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
                name:'resting',
                members:['alex', 'josh', 'lucas', 'peter'],
                dates:['24th august', '21th september'],
                description:'another description',
                checklists:[
                    {
                        id: 1,
                        name: 'checklist 1',
                        items:[
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
                        items:[
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
        id: 3,
        name: 'done',
        cards:[
            {   id: 1,
                name:'calculation',
                members:['alex', 'josh', 'lucas', 'peter'],
                dates:['24th august', '21th september'],
                description:'just some description',
                checklists:[
                    {
                        id: 1,
                        name: 'checklist 1',
                        items:[
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
                        items:[
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
                name:'checks',
                members:['alex', 'josh', 'lucas', 'peter'],
                dates:['24th august', '21th september'],
                description:'another description',
                checklists:[
                    {
                        id: 1,
                        name: 'checklist 1',
                        items:[
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
                        items:[
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
        id: 4,
        name: 'plan',
        cards:[
            {   id: 1,
                name:'meeting',
                members:['alex', 'josh', 'lucas', 'peter'],
                dates:['24th august', '21th september'],
                description:'just some description',
                checklists:[
                    {
                        id: 1,
                        name: 'checklist 1',
                        items:[
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
                        items:[
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
                name:'verdicts',
                members:['alex', 'josh', 'lucas', 'peter'],
                dates:['24th august', '21th september'],
                description:'another description',
                checklists:[
                    {
                        id: 1,
                        name: 'checklist 1',
                        items:[
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
                        items:[
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
]

export default lists