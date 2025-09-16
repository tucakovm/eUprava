export interface UUID extends String {}

export interface Dom {
  id: string;
  naziv: string;
  adresa: string;
}

export interface Student {
  id: string;
  ime: string;
  prezime: string;
  sobaId?: string | null;
}

export interface Soba {
  id: string;
  broj: string;
  slobodna: boolean;
  domId: string;
  kapacitet: number;
  studenti?: Student[];
  recenzije?: RecenzijaSobe[];
  kvarovi?: Kvar[];
}

export interface RecenzijaSobe {
  id: string;
  ocena: number;
  komentar?: string | null;
  sobaId: string;
  autorId: string;
}

export type StatusKvara = 'prijavljen' | 'u_toku' | 'resen';

export interface Kvar {
  id: string;
  opis: string;
  status: StatusKvara;
  sobaId: string;
  prijavioId: string;
}

export interface StudentskaKartica {
  id: string;
  stanje: number;
  studentID: string; // server šalje 'studentID' (camel case kao u domen modelu)
}