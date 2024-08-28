import { IncidentTimeline } from '@services/server';

const data: IncidentTimeline[] = [
  {
    change: {},
    createdAt: 1724775990000,
    id: 'incidentTimeline1',
    type: 'addComment',
    updatedAt: 1724775990000,
    userDetails: {
      email: 'admin@respondnow.io',
      name: 'Admin User',
      source: 'slack',
      userId: 'admin',
      userName: 'admin'
    }
  },
  {
    change: {},
    createdAt: 1724862390000,
    id: 'incidentTimeline2',
    type: 'addComment',
    updatedAt: 1724862390000,
    userDetails: {
      email: 'admin@respondnow.io',
      name: 'Admin User',
      source: 'slack',
      userId: 'admin',
      userName: 'admin'
    }
  }
];

export default data;
