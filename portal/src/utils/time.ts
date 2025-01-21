import moment from 'moment';
import { Incident } from '@services/server';

export function getDetailedTime(time: string | number, gethours?: boolean): string {
  return moment(time).format(gethours ? 'D MMM YYYY, HH:mm' : 'D MMM YYYY');
}

export function getDurationBasedOnStatus(startTime: number, endTime: number, status: Incident['status']): string {
  const currentTime = new Date().getTime();

  switch (status) {
    case 'Resolved':
      return moment.duration(endTime - startTime).humanize();
    default:
      return moment.duration(currentTime - startTime).humanize();
  }
}
