import { ComponentFixture, TestBed } from '@angular/core/testing';

import { Faults } from './faults';

describe('Faults', () => {
  let component: Faults;
  let fixture: ComponentFixture<Faults>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Faults]
    })
    .compileComponents();

    fixture = TestBed.createComponent(Faults);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
