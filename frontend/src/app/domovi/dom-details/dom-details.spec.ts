import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DomDetails } from './dom-details';

describe('DomDetails', () => {
  let component: DomDetails;
  let fixture: ComponentFixture<DomDetails>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DomDetails]
    })
    .compileComponents();

    fixture = TestBed.createComponent(DomDetails);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
