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
  username:string
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
  autorUsername: string;
}

export type StatusKvara = 'prijavljen' | 'u_toku' | 'resen';

export interface Kvar {
  id: string;
  opis: string;
  status: StatusKvara;
  sobaId: string;
  prijavioUsername: string;
}

export interface StudentskaKartica {
  id: string;
  stanje: number;
  studentUsername: string; 
}