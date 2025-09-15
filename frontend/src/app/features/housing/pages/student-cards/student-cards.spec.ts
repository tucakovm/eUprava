import { ComponentFixture, TestBed } from '@angular/core/testing';

import { StudentCards } from './student-cards';

describe('StudentCards', () => {
  let component: StudentCards;
  let fixture: ComponentFixture<StudentCards>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [StudentCards]
    })
    .compileComponents();

    fixture = TestBed.createComponent(StudentCards);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
