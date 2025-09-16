import { ComponentFixture, TestBed } from '@angular/core/testing';

import { Domovi } from './domovi';

describe('Domovi', () => {
  let component: Domovi;
  let fixture: ComponentFixture<Domovi>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Domovi]
    })
    .compileComponents();

    fixture = TestBed.createComponent(Domovi);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
