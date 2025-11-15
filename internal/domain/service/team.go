package service

import (
    "context"
    "fmt"
    "github.com/kimvlry/avito-internship-assignment/internal/domain"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/entity"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/repository"
)

type TeamService struct {
    teamRepository repository.TeamRepository
    userRepository repository.UserRepository
    tx             Transactor
}

func NewTeamService(teamRepo repository.TeamRepository, userRepo repository.UserRepository,
    tx Transactor) *TeamService {
    return &TeamService{
        teamRepository: teamRepo,
        userRepository: userRepo,
        tx:             tx,
    }
}

func (s *TeamService) CreateTeam(ctx context.Context, team *entity.Team, members []entity.User) (*entity.Team, error) {
    userIDs := make([]string, len(members))
    for i, member := range members {
        userIDs[i] = member.ID
    }

    err := s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
        exists, err := s.teamRepository.Exists(txCtx, team.Name)
        if err != nil {
            return fmt.Errorf("check team exists: %w", err)
        }
        if exists {
            return domain.ErrTeamAlreadyExists
        }

        if err := s.userRepository.CheckUsersAvailableForTeam(txCtx, userIDs, team.Name); err != nil {
            return err
        }

        for i := range members {
            member := &members[i]
            member.TeamName = team.Name

            exists, err := s.userRepository.Exists(txCtx, member.ID)
            if err != nil {
                return fmt.Errorf("check user exists: %w", err)
            }

            if exists {
                if err := s.userRepository.Update(txCtx, member); err != nil {
                    return fmt.Errorf("update user: %w", err)
                }
            } else {
                if err := s.userRepository.Create(txCtx, member); err != nil {
                    return fmt.Errorf("create user: %w", err)
                }
            }
        }

        if err := s.teamRepository.Create(txCtx, team); err != nil {
            return fmt.Errorf("create team: %w", err)
        }
        return nil
    })

    if err != nil {
        return nil, err
    }
    return team, nil
}

func (s *TeamService) GetTeamWithMembers(ctx context.Context, teamName string) (*entity.Team, []entity.User, error) {
    team, err := s.teamRepository.GetByName(ctx, teamName)
    if err != nil {
        return nil, nil, fmt.Errorf("get team: %w", err)
    }

    members, err := s.userRepository.GetByTeam(ctx, teamName)
    if err != nil {
        return nil, nil, fmt.Errorf("get users: %w", err)
    }
    return team, members, nil
}
