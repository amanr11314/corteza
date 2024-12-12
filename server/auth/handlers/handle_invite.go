package handlers

import (
	"github.com/cortezaproject/corteza/server/auth/request"
	"github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/cortezaproject/corteza/server/pkg/label"
	"go.uber.org/zap"
	"log"
)

func (h *AuthHandlers) acceptInviteForm(req *request.AuthReq) (err error) {
	h.Log.Debug("invite email password reset form")

	req.Template = TmplInvite

	// user not set, expecting valid token in URL
	if token := req.Request.URL.Query().Get("token"); len(token) > 0 {
		var user *types.User

		user, err = h.AuthService.ValidateInviteEmailToken(req.Context(), token)
		if err == nil {
			// login user
			req.AuthUser = request.NewAuthUser(h.Settings, user, false)

			if req.AuthUser.PendingEmailOTP() {
				// Email OTP enabled & pending
				//
				// If we're here it means user clicked on a link in an email;
				// we are effectively confirming email OTP
				req.AuthUser.CompleteEmailOTP()
			}

			// redirect back to self (but without token and with user in session)
			h.Log.Debug("valid user invite password reset token found, refreshing page with stored user")
			req.RedirectTo = GetLinks().AcceptInvite
			req.AuthUser.Save(req.Session)
			return nil
		}
	}

	if req.AuthUser == nil || err != nil {
		h.Log.Warn("invalid, user invite password reset token used", zap.Error(err))
		req.RedirectTo = GetLinks().Login
		t := translator(req, "auth")
		req.NewAlerts = append(req.NewAlerts, request.Alert{
			Type: "warning",
			Text: t("invite.alert.invalid-expired-invite-token"),
		})
	}

	req.Data["form"] = req.PopKV()
	return nil
}

func (h *AuthHandlers) acceptInviteProc(req *request.AuthReq) (err error) {
	h.Log.Debug("password reset proc")

	err = h.AuthService.SetPassword(req.Context(), req.AuthUser.User.ID, req.Request.PostFormValue("password"))

	if err == nil {
		t := translator(req, "auth")
		req.NewAlerts = append(req.NewAlerts, request.Alert{
			Type: "primary",
			Text: t("invite.alert.success"),
		})

		// delete all invite-token credential for this user
		credentialSvc := service.Credentials()
		cc, err := credentialSvc.FindUserInviteCrendential(req.Context(), req.AuthUser.User.ID)
		if err != nil {
			log.Printf("Error in Finding User Invite credentials: %v", err)
		} else {
			for _, userCredential := range cc {
				if err = h.Store.DeleteCredentialByID(req.Context(), userCredential.ID); err != nil {
					log.Printf("Error in Deleting User Invite credentials: %v", err)
				}
			}
		}

		// create inviteAccepted label for this user
		inviteAcceptedLabelresource := types.NewSimpleLabeledResource(req.AuthUser.User.ID, "user")


		inviteAcceptedLabelresource.SetLabel(types.InviteAcceptedLabel, "true")
		
		// update user label with key
		// "inviteAccepted": "true"
		if err = label.Create(req.Context(), h.Store, inviteAcceptedLabelresource); err != nil {
			log.Printf("Error in inviteAccepted label.Create: %v", err)
		}

		req.RedirectTo = GetLinks().Profile
		return nil
	}

	switch {
	case service.AuthErrPasswordResetDisabledByConfig().Is(err):
		h.passwordResetDisabledAlert(req)
		return nil

	default:
		h.Log.Error("unhandled error", zap.Error(err))
		return err
	}
}
