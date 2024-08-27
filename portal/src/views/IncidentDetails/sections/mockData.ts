import { IncidentTimeline } from '@services/server';

const data: IncidentTimeline[] = [
  {
    change: {},
    createdAt: 1700000000,
    id: 'incidentTimeline1',
    type: 'addComment',
    updatedAt: 1700000000,
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
