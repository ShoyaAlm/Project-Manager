const lists = [
    {
        id: 1,
        name: 'درحال انجام',
        cards:[
            {   id: 1,
                name:'برنامه ریزی',
                members:['پرهام', 'شایان', 'محمد', 'اشکان'],
                dates:['30 مرداد', '30 شهریور'],
                description:'یک سری توضیحات',
                checklists:[
                    {
                        id: 1,
                        name: 'چکلیست 1',
                        items:[
                            {
                                id: 1,
                                name:'آیتم 1',
                                dueDate:'20 شهریور',
                                assignedTo:['محمد', 'اشکان']
                            },
                            {
                                id: 2,
                                name:'آیتم 2',
                                dueDate:'25 شهریور',
                                assignedTo:['پرهام', 'شایان']
                            },

                        ]
                    },
                    {
                        id: 2,
                        name: 'چکلیست 2',
                        items:[
                            {
                                id: 1,
                                name:'آیتم 1',
                                dueDate:'20 شهریور',
                                assignedTo:['محمد', 'اشکان']
                            },
                            {
                                id: 2,
                                name:'آیتم 2',
                                dueDate:'25 شهریور',
                                assignedTo:['پرهام', 'شایان']
                            },

                        ]
                    }
                ]
            },
            {   id: 2,
                name:'نوشتن برنامه',
                members:['پرهام', 'شایان', 'محمد', 'اشکان'],
                dates:['30 مرداد', '30 شهریور'],
                description:'یک سری توضیحات',
                checklists:[
                    {
                        id: 1,
                        name: 'چکلیست 1',
                        items:[
                            {
                                id: 1,
                                name:'آیتم 1',
                                dueDate:'20 شهریور',
                                assignedTo:['محمد', 'اشکان']
                            },
                            {
                                id: 2,
                                name:'آیتم 2',
                                dueDate:'25 شهریور',
                                assignedTo:['پرهام', 'شایان']
                            },

                        ]
                    },
                    {
                        id: 2,
                        name: 'چکلیست 2',
                        items:[
                            {
                                id: 1,
                                name:'آیتم 1',
                                dueDate:'20 شهریور',
                                assignedTo:['محمد', 'اشکان']
                            },
                            {
                                id: 2,
                                name:'آیتم 2',
                                dueDate:'25 شهریور',
                                assignedTo:['پرهام', 'شایان']
                            },

                        ]
                    }
                ]
            },
        ]
    },
    {
        id: 1,
        name: 'انجام شده',
        cards:[
            {   id: 1,
                name:'برنامه ریزی',
                members:['پرهام', 'شایان', 'محمد', 'اشکان'],
                dates:['30 مرداد', '30 شهریور'],
                description:'یک سری توضیحات',
                checklists:[
                    {
                        id: 1,
                        name: 'چکلیست 1',
                        items:[
                            {
                                id: 1,
                                name:'آیتم 1',
                                dueDate:'20 شهریور',
                                assignedTo:['محمد', 'اشکان']
                            },
                            {
                                id: 2,
                                name:'آیتم 2',
                                dueDate:'25 شهریور',
                                assignedTo:['پرهام', 'شایان']
                            },

                        ]
                    },
                    {
                        id: 2,
                        name: 'چکلیست 2',
                        items:[
                            {
                                id: 1,
                                name:'آیتم 1',
                                dueDate:'20 شهریور',
                                assignedTo:['محمد', 'اشکان']
                            },
                            {
                                id: 2,
                                name:'آیتم 2',
                                dueDate:'25 شهریور',
                                assignedTo:['پرهام', 'شایان']
                            },

                        ]
                    }
                ]
            },
            {   id: 2,
                name:'نوشتن برنامه',
                members:['پرهام', 'شایان', 'محمد', 'اشکان'],
                dates:['30 مرداد', '30 شهریور'],
                description:'یک سری توضیحات',
                checklists:[
                    {
                        id: 1,
                        name: 'چکلیست 1',
                        items:[
                            {
                                id: 1,
                                name:'آیتم 1',
                                dueDate:'20 شهریور',
                                assignedTo:['محمد', 'اشکان']
                            },
                            {
                                id: 2,
                                name:'آیتم 2',
                                dueDate:'25 شهریور',
                                assignedTo:['پرهام', 'شایان']
                            },

                        ]
                    },
                    {
                        id: 2,
                        name: 'چکلیست 2',
                        items:[
                            {
                                id: 1,
                                name:'آیتم 1',
                                dueDate:'20 شهریور',
                                assignedTo:['محمد', 'اشکان']
                            },
                            {
                                id: 2,
                                name:'آیتم 2',
                                dueDate:'25 شهریور',
                                assignedTo:['پرهام', 'شایان']
                            },

                        ]
                    }
                ]
            },
        ]
    },
    {
        id: 1,
        name: 'به پایان رسیده',
        cards:[
            {   id: 1,
                name:'برنامه ریزی',
                members:['پرهام', 'شایان', 'محمد', 'اشکان'],
                dates:['30 مرداد', '30 شهریور'],
                description:'یک سری توضیحات',
                checklists:[
                    {
                        id: 1,
                        name: 'چکلیست 1',
                        items:[
                            {
                                id: 1,
                                name:'آیتم 1',
                                dueDate:'20 شهریور',
                                assignedTo:['محمد', 'اشکان']
                            },
                            {
                                id: 2,
                                name:'آیتم 2',
                                dueDate:'25 شهریور',
                                assignedTo:['پرهام', 'شایان']
                            },

                        ]
                    },
                    {
                        id: 2,
                        name: 'چکلیست 2',
                        items:[
                            {
                                id: 1,
                                name:'آیتم 1',
                                dueDate:'20 شهریور',
                                assignedTo:['محمد', 'اشکان']
                            },
                            {
                                id: 2,
                                name:'آیتم 2',
                                dueDate:'25 شهریور',
                                assignedTo:['پرهام', 'شایان']
                            },

                        ]
                    }
                ]
            },
            {   id: 2,
                name:'نوشتن برنامه',
                members:['پرهام', 'شایان', 'محمد', 'اشکان'],
                dates:['30 مرداد', '30 شهریور'],
                description:'یک سری توضیحات',
                checklists:[
                    {
                        id: 1,
                        name: 'چکلیست 1',
                        items:[
                            {
                                id: 1,
                                name:'آیتم 1',
                                dueDate:'20 شهریور',
                                assignedTo:['محمد', 'اشکان']
                            },
                            {
                                id: 2,
                                name:'آیتم 2',
                                dueDate:'25 شهریور',
                                assignedTo:['پرهام', 'شایان']
                            },

                        ]
                    },
                    {
                        id: 2,
                        name: 'چکلیست 2',
                        items:[
                            {
                                id: 1,
                                name:'آیتم 1',
                                dueDate:'20 شهریور',
                                assignedTo:['محمد', 'اشکان']
                            },
                            {
                                id: 2,
                                name:'آیتم 2',
                                dueDate:'25 شهریور',
                                assignedTo:['پرهام', 'شایان']
                            },

                        ]
                    }
                ]
            },
        ]
    },
    {
        id: 1,
        name: 'پروژه ها',
        cards:[
            {   id: 1,
                name:'برنامه ریزی',
                members:['پرهام', 'شایان', 'محمد', 'اشکان'],
                dates:['30 مرداد', '30 شهریور'],
                description:'یک سری توضیحات',
                checklists:[
                    {
                        id: 1,
                        name: 'چکلیست 1',
                        items:[
                            {
                                id: 1,
                                name:'آیتم 1',
                                dueDate:'20 شهریور',
                                assignedTo:['محمد', 'اشکان']
                            },
                            {
                                id: 2,
                                name:'آیتم 2',
                                dueDate:'25 شهریور',
                                assignedTo:['پرهام', 'شایان']
                            },

                        ]
                    },
                    {
                        id: 2,
                        name: 'چکلیست 2',
                        items:[
                            {
                                id: 1,
                                name:'آیتم 1',
                                dueDate:'20 شهریور',
                                assignedTo:['محمد', 'اشکان']
                            },
                            {
                                id: 2,
                                name:'آیتم 2',
                                dueDate:'25 شهریور',
                                assignedTo:['پرهام', 'شایان']
                            },

                        ]
                    }
                ]
            },
            {   id: 2,
                name:'نوشتن برنامه',
                members:['پرهام', 'شایان', 'محمد', 'اشکان'],
                dates:['30 مرداد', '30 شهریور'],
                description:'یک سری توضیحات',
                checklists:[
                    {
                        id: 1,
                        name: 'چکلیست 1',
                        items:[
                            {
                                id: 1,
                                name:'آیتم 1',
                                dueDate:'20 شهریور',
                                assignedTo:['محمد', 'اشکان']
                            },
                            {
                                id: 2,
                                name:'آیتم 2',
                                dueDate:'25 شهریور',
                                assignedTo:['پرهام', 'شایان']
                            },

                        ]
                    },
                    {
                        id: 2,
                        name: 'چکلیست 2',
                        items:[
                            {
                                id: 1,
                                name:'آیتم 1',
                                dueDate:'20 شهریور',
                                assignedTo:['محمد', 'اشکان']
                            },
                            {
                                id: 2,
                                name:'آیتم 2',
                                dueDate:'25 شهریور',
                                assignedTo:['پرهام', 'شایان']
                            },

                        ]
                    }
                ]
            },
        ]
    },
]

export default lists